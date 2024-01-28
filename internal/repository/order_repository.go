package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

var OrderNotFoundErr = errors.New("order not found")
var AdditionsNotFoundErr = errors.New("some of addition inforamation about order not found")

type OrderRepository struct {
	psqlDB  *sqlx.DB
	redisDB *redis.Client
}

func NewOrderRepository(psqlDB *sqlx.DB, redisDB *redis.Client) *OrderRepository {
	return &OrderRepository{psqlDB: psqlDB, redisDB: redisDB}
}
