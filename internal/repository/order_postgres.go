package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"nats_server/internal/entity"
)

var OrderNotFoundErr = errors.New("order not found")
var AdditionsNotFoundErr = errors.New("some of addition inforamation about order not found")

type OrderRepository struct {
	psqlDb *sqlx.DB
}

func NewOrderPostgres(psqlDb *sqlx.DB) *OrderRepository {
	return &OrderRepository{psqlDb: psqlDb}
}

func (r *OrderRepository) Create(order entity.Order) (string, error) {
	return "", nil
}

func (r *OrderRepository) GetWithAddition(orderUID string) (entity.Order, error) {
	var order entity.Order

	err := r.psqlDb.Get(&order, "SELECT * FROM orders WHERE order_uid = $1", orderUID)
	if err != nil {
		return order, OrderNotFoundErr
	}

	err = r.enrichOrder(&order)

	return order, err
}

func (r *OrderRepository) GetWithAdditionList() ([]entity.Order, error) {
	var list []entity.Order

	err := r.psqlDb.Select(&list, "SELECT * FROM orders")
	if err != nil {
		return nil, err
	}

	for i := range list {
		err = r.enrichOrder(&list[i])
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

// enrichOrder enriches order with addition information
func (r *OrderRepository) enrichOrder(order *entity.Order) error {
	err := r.psqlDb.Get(&order.Delivery, "SELECT * FROM delivery WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		return AdditionsNotFoundErr
	}

	err = r.psqlDb.Get(&order.Payment, "SELECT * FROM payment WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		return AdditionsNotFoundErr
	}

	err = r.psqlDb.Select(&order.Items, "SELECT * FROM items WHERE order_uid = $1", order.OrderUID)

	return nil
}
