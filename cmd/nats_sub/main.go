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
	pkgRepo "nats_server/pkg/repository"
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

	psqlDb, err := pkgRepo.NewPostgresDB(pkgRepo.PSQLConfig{
		Port:     os.Getenv("PSQL_PORT"),
		Username: os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   os.Getenv("PSQL_DB"),
		SSLMode:  "disable",
	})
	if err != nil {
		panic(err)
	}

	repos := repository.NewRepository(psqlDb)
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
			// Do not unsubscribe a durable on exit, except if asked to.
			//if durable == "" || unsubscribe {
			//	sub.Unsubscribe()
			//}
			sub.Unsubscribe()
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
