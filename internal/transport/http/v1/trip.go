package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/itoqsky/InnoCoTravel-backend/docs"
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/internal/server"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/response"
)

func (h *Handler) initTripsRoutes(api *gin.RouterGroup) {
	trip := api.Group("/trip", h.userIdentity)
	{
		trip.POST("/", h.createTrip)
		trip.DELETE("/:id", h.deleteTrip)

		trip.GET("/:id", h.getTrip)
		trip.GET("/", h.getJoinedTrips)
		trip.GET("/:id/users", h.getJoinedUsers)

		trip.PUT("/adjacent", h.getAdjacentTrips)
	}
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) createTrip(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var trip core.Trip
	if err := c.BindJSON(&trip); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trip.AdminId = uctx.UserId
	trip.AdminUsername = uctx.Username

	tripId, err := h.services.Trip.Create(trip)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.hub.Rooms[tripId] = &server.Room{
		Id:      tripId,
		Name:    getTripName(trip.FromPoint, trip.ToPoint, trip.ChosenTimestamp),
		Clients: make(map[int64]*server.Client),
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"trip_id": tripId,
	})
}

func getNameOfPoint(p int) string {
	if p == 1 {
		return "INNO"
	} else if p == 2 {
		return "KZN"
	} else if p == 3 {
		return "AIRPORT_KZN"
	}
	return "BRUH"
}

func getTripName(from, to int, timestamp string) string {
	return fmt.Sprintf("%s -> %s at:%s", getNameOfPoint(from), getNameOfPoint(to), timestamp)
}

func (h *Handler) getJoinedTrips(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	trips, err := h.services.Trip.GetJoinedTrips(uctx.UserId)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: trips})
}

func (h *Handler) getTrip(c *gin.Context) {
	// uctx, err := getUserCtx(c)
	// if err != nil {
	// 	response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	tripId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trip, err := h.services.Trip.GetById(int64(tripId))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, trip)
}

func (h *Handler) deleteTrip(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tripId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newAdminId, err := h.services.Trip.Delete(uctx.UserId, int64(tripId))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"new_admin_id": newAdminId,
	})
}

func (h *Handler) getAdjacentTrips(c *gin.Context) {
	var input core.InputAdjTrips
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trips, err := h.services.Trip.GetAdjTrips(input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: trips})
}

func (h *Handler) getJoinedUsers(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tripId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.services.Trip.GetJoinedUsers(uctx.UserId, int64(tripId))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: users})
}
