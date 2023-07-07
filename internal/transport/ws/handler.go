package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel-backend/internal/server"
)

type Handler struct {
	hub *server.Hub
}

func NewHandler(hub *server.Hub) *Handler {
	return &Handler{hub}
}

func (h *Handler) InitRoutes() {
	r := gin.Default()

	ws := r.Group("/ws")
	{
		ws.POST("/createRoom", h.createRoom)
		// ws.GET("/joinRoom/:roomId", joinRoom)
		// ws.GET("/getRooms", getRooms)
		// ws.GET("/getClients/:roomId", getClients)
	}
}
