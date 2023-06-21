package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel_backend/internal/core"
)

func (h *Handler) createTrip(c *gin.Context) {
	userId := getUserId(c)

	var trip core.Trip
	if err := c.BindJSON(&trip); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trip.AdminId = userId

	tripId, err := h.services.Trip.Create(trip)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"trip_id": tripId,
	})
}

func (h *Handler) getTrip(c *gin.Context) {
	userId := getUserId(c)
	tripId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trip, err := h.services.Trip.GetById(userId, tripId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, trip)
}

func (h *Handler) deleteTrip(c *gin.Context) {
	userId := getUserId(c)
	tripId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newAdminId, err := h.services.Trip.Delete(userId, tripId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"new_admin_id": newAdminId,
	})
}

// func (h *Handler) updateTrip(c *gin.Context) {
// 	userId := getUserId(c)
// 	tripId, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	var trip core.Trip
// 	if err := c.BindJSON(&trip); err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	trip.AdminId = userId
// 	trip.TripId = tripId
// 	if err := h.services.Trip.Update(trip); err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, statusOkResponse{"ok"})
// }

type resAdjTrips struct {
	Data []core.Trip `json:"data"`
}

func (h *Handler) getAdjacentTrips(c *gin.Context) {
	// userId := getUserId(c)

	var input core.InputAdjTrips
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trips, err := h.services.Trip.GetAdjTrips(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resAdjTrips{trips})
}
