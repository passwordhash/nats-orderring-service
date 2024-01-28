package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"nats_server/internal/nats"
	"nats_server/internal/repository"
	"nats_server/internal/service"
	"os"
	"os/signal"
)

var URL = "nats://localhost:4222"
var clusterID = "test-cluster"
var clientID = "test-client"

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	psqlDb, err := repository.NewDefaultPsqlDB()
	if err != nil {
		panic(err)
	}
	redisDb, err := repository.NewDefaultRedisDB()
	if err != nil {
		panic(err)
	}

	repos := repository.NewRepository(psqlDb, redisDb)
	services := service.NewService(repos)

	clientID := uuid.New().String()

	// Connect to Nats
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	defer sc.Close()

	log.Println("Connected to " + URL)

	// Subscribe to subject
	startOpt := stan.StartAt(0)

	sub, err := sc.QueueSubscribe("foo", "bar", func(msg *stan.Msg) {
		nats.HandleMessage(msg, services)
	}, startOpt)
	if err != nil {
		sc.Close()
		log.Fatalf("Error during subscribe: %v\n", err)
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			sub.Unsubscribe()
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
