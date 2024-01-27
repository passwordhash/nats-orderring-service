package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"nats_server/internal/handler"
	"nats_server/pkg/server"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("cannot load env file: %v", err.Error())
	}

	config := NewConfig()

	handlers := handler.NewHandler()

	srv := new(server.Server)

	logrus.Infof("Try to start HTTP Server on port %s ...", config.HttpPort)
	if err := srv.Run(config.HttpPort, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

}

type Config struct {
	HttpPort string
}

func NewConfig() *Config {
	return &Config{
		HttpPort: os.Getenv("HTTP_PORT"),
	}
}
