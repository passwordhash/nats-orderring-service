package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"log"
	"nats_server/internal/entity"
	"os"
)

var URL = "nats://localhost:4222"
var clusterID = "test-cluster"

func main() {
	var order entity.Order

	b, err := os.ReadFile("./cmd/nats_pub/input.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(b, &order); err != nil {
		panic(err)
	}

	clientID := uuid.New().String()

	// Connect to Nats
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	defer sc.Close()

	log.Println("Connected to " + URL)

	if err := sc.Publish("foo", []byte("heeasfasf")); err != nil {
		log.Printf("Error during publish: %v\n", err)
	}
}
