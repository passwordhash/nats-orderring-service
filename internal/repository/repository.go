package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"nats_server/internal/entity"
	pkgRepo "nats_server/pkg/repository"
	"os"
)

type Order interface {
	Create(order entity.Order) (string, error)

	GetWithAddition(orderUID string) (entity.Order, error)
	GetWithAdditionList() ([]entity.Order, error)

	Save(order entity.Order) error
	Delete(orderUID string) error
	Read(orderUID string) (entity.Order, error)
}

type Repository struct {
	Order *OrderRepository
}

func NewRepository(psqlDB *sqlx.DB, redisDB *redis.Client) *Repository {
	return &Repository{
		Order: NewOrderRepository(psqlDB, redisDB),
	}
}

func NewDefaultPsqlDB() (*sqlx.DB, error) {
	return pkgRepo.NewPostgresDB(pkgRepo.PSQLConfig{
		Port:     os.Getenv("PSQL_PORT"),
		Username: os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   os.Getenv("PSQL_DB"),
		SSLMode:  "disable",
	})
}

func NewDefaultRedisDB() (*redis.Client, error) {
	addr := "localhost:" + os.Getenv("REDIS_PORT")
	db := pkgRepo.NewRedisDB(pkgRepo.RedisConfig{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if _, err := db.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return db, nil
}
