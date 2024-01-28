package repository

import (
	"context"
	"encoding/json"
	"nats_server/internal/entity"
	"time"
)

var cacheTTL = 4 * time.Hour

func (r *OrderRepository) Save(order entity.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = r.redisDB.Set(context.TODO(), order.OrderUID, orderJSON, cacheTTL).Err()

	return err
}

func (r *OrderRepository) Delete(orderUID string) error {
	return nil
}

func (r *OrderRepository) Read(orderUID string) (entity.Order, error) {
	var order entity.Order

	orderJSON, err := r.redisDB.Get(context.TODO(), orderUID).Result()
	if err != nil {
		return order, err
	}

	if err := json.Unmarshal([]byte(orderJSON), &order); err != nil {
		return entity.Order{}, err
	}

	return order, nil
}
