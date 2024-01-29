package main

import (
	"encoding/json"
	"flag"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"log"
	"nats_server/internal/entity"
	"os"
)

var (
	defaultURL = "nats://localhost:4222"
)
var clusterID = "test-cluster"

func main() {

	var (
		inputBaseDir string
		natsURL      string
	)

	flag.StringVar(&natsURL, "url", defaultURL, "Nats URL")
	flag.StringVar(&inputBaseDir, "dir", "", "Input files base dir")
	flag.Parse()

	posArgs := 2
	if natsURL != defaultURL {
		posArgs++
	}

	if posArgs >= len(os.Args) {
		logrus.Fatal("Input files not found")
	}

	if inputBaseDir == "" {
		logrus.Fatal("Input files base dir is empty")
	}

	clientID := uuid.New().String()

	// Connect to Nats
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsURL)
	}
	defer sc.Close()

	log.Println("Connected to " + natsURL)

	for _, v := range os.Args[posArgs:] {
		b, err := readInput("./cmd/nats_pub/mocks/" + v)
		if err != nil {
			logrus.Error("Error during read input")
			continue
		}

		if err := sc.Publish("foo", b); err != nil {
			log.Printf("Error during publish: %v\n", err)
		}

		logrus.Info("Message published")
	}

}

func readInput(path string) ([]byte, error) {
	var order entity.Order

	b, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	if err := json.Unmarshal(b, &order); err != nil {
		logrus.Error("Error during unmarshal")
		return []byte{}, err
	}

	return b, nil
}
