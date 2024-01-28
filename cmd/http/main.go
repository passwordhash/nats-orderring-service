package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"nats_server/internal/handler"
	"nats_server/internal/repository"
	"nats_server/internal/service"
	"nats_server/pkg/server"
	"os"
)

func main() {
	// Load env variables
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("cannot load env file: %v", err.Error())
	}

	config := NewConfig()

	// Connect to PostgreSQL
	psqlDB, err := repository.NewDefaultPsqlDB()
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer psqlDB.Close()
	logrus.Info("PostgreSQL connected")

	// Connect to Redis
	redisDB, err := repository.NewDefaultRedisDB()
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer redisDB.Close()
	logrus.Info("Redis connected")

	// Initialize repositories, services and handlers
	repos := repository.NewRepository(psqlDB, redisDB)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	// Initialize HTTP Server
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
