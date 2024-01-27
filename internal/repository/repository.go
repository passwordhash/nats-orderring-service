package repository

import (
	"github.com/jmoiron/sqlx"
	"nats_server/internal/entity"
)

type Order interface {
	Create(order entity.Order) (string, error)

	GetWithAddition(orderUID string) (entity.Order, error)
	GetWithAdditionList() ([]entity.Order, error)
}

type Repository struct {
	Order *OrderRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order: NewOrderPostgres(db),
	}
}
