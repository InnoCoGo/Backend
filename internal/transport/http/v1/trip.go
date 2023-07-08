package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/itoqsky/InnoCoTravel-backend/docs"
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
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

// createTrip 	godoc
//
//	@Summary		Create trip
//	@Tags			trips
//	@Description	create trip
//	@Security		ApiKeyAuth
//	@ID				create-trip
//	@Accept			json
//	@Produce		json
//	@Param			createInput	body		core.Trip	true	"trip info"
//	@Success		200			{integer}	integer
//	@Failure		400			{object}	errorResponse
//	@Failure		404			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Failure		default		{object}	errorResponse
//	@Router			/trips [post]

func (h *Handler) createTrip(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var trip core.Trip
	if err := c.BindJSON(&trip); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trip.AdminId = uctx.UserId
	trip.AdminUsername = uctx.Username

	tripId, err := h.services.Trip.Create(trip)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"trip_id": tripId,
	})
}

// 	getJoinedTrips 			godoc
//	@Summary		Get Joined Trips
//	@Tags			trips
//	@Description	get all trips
//	@ID				getjoinedTrips
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			getJoinedTrips	body		int	false	"bruh"
//	@Success		200				{object}	dataResponse
//	@Failure		400				{object}	errorResponse
//	@Failure		404				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Failure		default			{object}	errorResponse
//	@Router			/trips [get]

func (h *Handler) getJoinedTrips(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	trips, err := h.services.Trip.GetJoinedTrips(uctx.UserId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataResponse{Data: trips})
}

// 	getTrip 			godoc
//	@Summary		Get Trip
//	@Tags			trips
//	@Description	get trip
//	@Security		ApiKeyAuth
//	@ID				get-trip
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"Trip ID"
//	@Success		200		{object}	core.Trip
//	@Failure		400		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Failure		default	{object}	errorResponse
//	@Router			/trips/{id} [get]

func (h *Handler) getTrip(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tripId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trip, err := h.services.Trip.GetById(uctx.UserId, tripId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, trip)
}

// 	deleteTrip 			godoc
//	@Summary		Delete Trip
//	@Tags			trips
//	@Description	delete trip
//	@Security		ApiKeyAuth
//	@ID				delete-trip
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"Trip ID"
//	@Success		200		{object}	integer
//	@Failure		400		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Failure		default	{object}	errorResponse
//	@Router			/trips/{id} [delete]

func (h *Handler) deleteTrip(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tripId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newAdminId, err := h.services.Trip.Delete(uctx.UserId, tripId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"new_admin_id": newAdminId,
	})
}

// 	getAdjacentTrips 		godoc
//	@Summary		Get Adjacent Trips
//	@Tags			trips
//	@Description	Get Adjacent Trips
//	@Security		ApiKeyAuth
//	@ID				get-adjacent-trips
//	@Accept			json
//	@Produce		json
//	@Param			getAdjacentTrips	body		core.InputAdjTrips	true	"adj trip info"
//	@Success		200					{object}	dataResponse
//	@Failure		400					{object}	errorResponse
//	@Failure		404					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Failure		default				{object}	errorResponse
//	@Router			/trips/adjacent [put]

func (h *Handler) getAdjacentTrips(c *gin.Context) {
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

	c.JSON(http.StatusOK, dataResponse{Data: trips})
}

func (h *Handler) getJoinedUsers(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tripId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.services.Trip.GetJoinedUsers(uctx.UserId, tripId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataResponse{Data: users})
}
