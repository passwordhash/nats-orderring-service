package handler

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	//services *service.Service
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) InitRoutes() *gin.Engine {
	app := gin.New()

	api := app.Group("/api")

	api.GET("/:id", h.get)

	api.GET("/", func(context *gin.Context) {

	})

	return app
}

func (h *Handler) get(c *gin.Context) {
	id := c.Param("id")

	c.JSON(200, gin.H{
		"id": id,
	})
}

func (h *Handler) getAll(c *gin.Context) {
}
