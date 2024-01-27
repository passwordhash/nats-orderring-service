package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"nats_server/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	app := gin.New()

	api := app.Group("/api")

	api.GET("/:id", h.get)

	api.GET("/", h.getAll)

	return app
}

func (h *Handler) get(c *gin.Context) {
	id := c.Param("id")

	order, err := h.services.Order.Get(id)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, order)
}

func (h *Handler) getAll(c *gin.Context) {
	orders, err := h.services.Order.GetList()
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, orders)
}
