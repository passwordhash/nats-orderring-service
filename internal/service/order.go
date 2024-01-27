package service

import (
	"nats_server/internal/entity"
	"nats_server/internal/repository"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) Get(orderUID string) (entity.Order, error) {
	return s.repo.GetWithAddition(orderUID)
}

func (s *OrderService) GetList() ([]entity.Order, error) {
	return s.repo.GetWithAdditionList()
}
