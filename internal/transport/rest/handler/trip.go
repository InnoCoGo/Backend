package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel_backend/internal/core"
)

func (h *Handler) createTrip(c *gin.Context) {
	userId, _ := getUserId(c)

	var trip core.Trip
	if err := c.BindJSON(&trip); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	trip.AdminId = userId

	tripId, err := h.services.Trip.Create(trip)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "empty user id or username")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"tripId": tripId,
	})

}
