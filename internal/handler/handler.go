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

	app.GET("/check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return app
}
