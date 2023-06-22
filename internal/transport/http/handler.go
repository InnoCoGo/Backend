package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel_backend/internal/service"
	v1 "github.com/itoqsky/InnoCoTravel_backend/internal/transport/http/v1"
)

type Handler struct {
	services *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{services: s}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		// TODO: Limiter
		corsMiddleware,
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)
	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services)
	api := router.Group("/api")
	{
		handlerV1.InitV1(api)
	}
}
