package service

import (
	"errors"
	"nats_server/internal/entity"
	"nats_server/internal/repository"
)

var OrderAlreadyExistsErr = errors.New("order already exists")

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) Create(order entity.Order) (string, error) {
	condidate, _ := s.repo.GetWithAddition(order.OrderUID)
	if condidate.OrderUID != "" {
		return "", OrderAlreadyExistsErr
	}
	return s.repo.Create(order)
}

func (s *OrderService) Get(orderUID string) (entity.Order, error) {
	return s.repo.GetWithAddition(orderUID)
}

func (s *OrderService) GetList() ([]entity.Order, error) {
	return s.repo.GetWithAdditionList()
}
