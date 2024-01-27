package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"nats_server/internal/entity"
	"nats_server/internal/repository"
	"nats_server/internal/service"
	pkgRepo "nats_server/pkg/repository"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	var order entity.Order

	input, err := os.ReadFile("./cmd/nats/input.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(input, &order); err != nil {
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

	if _, err := services.Create(order); err != nil {
		panic(err)
	}
}
