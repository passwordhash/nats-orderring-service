package service

import (
	"nats_server/internal/entity"
	"nats_server/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Order interface {
	Create(order entity.Order) (string, error)
	Get(orderUID string) (entity.Order, error)
	GetList() ([]entity.Order, error)
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repos.Order),
	}
}
