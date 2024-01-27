package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"nats_server/internal/handler"
	"nats_server/internal/repository"
	"nats_server/internal/service"
	pkgRepo "nats_server/pkg/repository"
	"nats_server/pkg/server"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("cannot load env file: %v", err.Error())
	}

	config := NewConfig()

	psqlDb, err := pkgRepo.NewPostgresDB(pkgRepo.PSQLConfig{
		Port:     os.Getenv("PSQL_PORT"),
		Username: os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   os.Getenv("PSQL_DB"),
		SSLMode:  "disable",
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	logrus.Info("PostgreSQL connected")

	repos := repository.NewRepository(psqlDb)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

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
