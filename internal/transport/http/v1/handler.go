package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel-backend/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{services: s}
}

func (h *Handler) InitV1(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initAuthRoutes(v1)
		// h.initUsersRoutes(v1)
		h.initTripsRoutes(v1)
	}
}
