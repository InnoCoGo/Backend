package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel-backend/internal/server"
	"github.com/itoqsky/InnoCoTravel-backend/internal/service"
)

type Handler struct {
	hub         *server.Hub
	authService *service.Service
}

func NewHandler(h *server.Hub, s *service.Service) *Handler {
	return &Handler{h, s}
}

func (h *Handler) InitWsRoutes() *gin.Engine {
	r := gin.Default()

	ws := r.Group("/ws")
	{
		ws.POST("", h.createRoom)
		ws.GET("", h.joinRoom)
		// ws.GET("/getRooms", getRooms)
		// ws.GET("/getClients/:roomId", getClients)
	}
	return r
}
