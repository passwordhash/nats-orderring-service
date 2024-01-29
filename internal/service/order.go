package service

import (
	"errors"
	"github.com/sirupsen/logrus"
	"nats_server/internal/entity"
	"nats_server/internal/repository"
)

var OrderAlreadyExistsErr = errors.New("order already exists")
var OrderNotFoundErr = errors.New("order not found")

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
	if !condidate.IsEmpty() {
		return "", OrderAlreadyExistsErr
	}

	uid, err := s.repo.Create(order)
	if err != nil {
		return "", err
	}

	return uid, s.repo.Save(order)
}

func (s *OrderService) Get(orderUID string) (entity.Order, error) {
	var order entity.Order

	// Try to get order from cache
	order, err := s.repo.Read(orderUID)
	if !order.IsEmpty() && err == nil {
		logrus.Info("Order found in cache ", order.OrderUID)
		return order, nil
	}

	order, err = s.repo.GetWithAddition(orderUID)
	if order.IsEmpty() {
		return order, errors.Join(OrderNotFoundErr, err)
	}
	if err != nil {
		return order, err
	}

	// Save order to cache
	if err := s.repo.Save(order); err != nil {
		logrus.Error("Error while saving order to cache ", err)
	}

	return order, nil
}

func (s *OrderService) GetList() ([]entity.Order, error) {
	return s.repo.GetWithAdditionList()
}
