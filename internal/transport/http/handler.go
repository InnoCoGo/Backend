package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/itoqsky/InnoCoTravel-backend/docs"
	"github.com/itoqsky/InnoCoTravel-backend/internal/server"
	"github.com/itoqsky/InnoCoTravel-backend/internal/service"
	v1 "github.com/itoqsky/InnoCoTravel-backend/internal/transport/http/v1"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/limiter"
)

type Handler struct {
	hub      *server.Hub
	services *service.Service
}

func NewHandler(s *service.Service, h *server.Hub) *Handler {
	return &Handler{services: s, hub: h}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		limiter.Limit(10, 20, 10*time.Minute),
		corsMiddleware,
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)
	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.hub)
	api := router.Group("/api")
	{
		handlerV1.InitV1(api)
	}

	// ws := router.Group("/ws")
	// {
	// 	handlerV1.InitV1(ws)
	// }
}
