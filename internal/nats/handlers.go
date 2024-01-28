package nats

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"nats_server/internal/entity"
	"nats_server/internal/service"
)

func HandleMessage(m *stan.Msg, service *service.Service) {
	var order entity.Order

	if err := json.Unmarshal(m.Data, &order); err != nil {
		logrus.Error(m)
		return
	}

	logrus.Info(order)

	if _, err := service.Create(order); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("New msg was saved!")
}
