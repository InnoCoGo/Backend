package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel_backend/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{services: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	routes := gin.New()

	auth := routes.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)

		auth.POST("/tg-login", h.tgLogIn)
	}

	api := routes.Group("/api", h.userIdentity)
	{
		trip := api.Group("/trip")
		{
			trip.POST("/", h.createTrip)
			trip.GET("/:id", h.getTrip)
			trip.DELETE("/:id", h.deleteTrip)
			trip.GET("/", h.getJoinedTrips)
			// trip.POST("/join", h.joinTrip)
			trip.GET("/adjacent", h.getAdjacentTrips)
		}
	}
	return routes
}
