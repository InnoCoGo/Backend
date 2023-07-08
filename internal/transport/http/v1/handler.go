package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel-backend/internal/server"
	"github.com/itoqsky/InnoCoTravel-backend/internal/service"
)

type Handler struct {
	hub      *server.Hub
	services *service.Service
}

func NewHandler(s *service.Service, h *server.Hub) *Handler {
	return &Handler{services: s, hub: h}
}

func (h *Handler) InitV1(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initAuthRoutes(v1)
		// h.initUsersRoutes(v1)
		h.initWsTripsRoutes(v1)
		h.initTripsRoutes(v1)
	}
}
