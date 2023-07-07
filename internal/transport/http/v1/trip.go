package v1

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/itoqsky/InnoCoTravel-backend/docs"
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/response"
)

func (h *Handler) initTripsRoutes(api *gin.RouterGroup) {
	trip := api.Group("/trip", h.userIdentity)
	{
		trip.POST("/", h.createTrip)
		trip.GET("/", h.getJoinedTrips)
		trip.GET("/:id", h.getTrip)
		trip.DELETE("/:id", h.deleteTrip)

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

	reqData := CreateRoomReq{
		ID:   strconv.Itoa(int(tripId)),
		Name: fmt.Sprintf("%s_to_%s_at_%s", getNameOfPoint(trip.FromPoint), getNameOfPoint(trip.ToPoint), trip.ChosenTimestamp),
	}

	if err := doRequest(reqData); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"trip_id": tripId,
	})
}

func doRequest(reqData CreateRoomReq) error {
	u := url.URL{
		Scheme: "http",
		Host:   os.Getenv("CHAT_HOST"),
		Path:   path.Join("ws", "createRoom"),
	}

	q := u.Query()
	q.Set("id", reqData.ID)
	q.Set("name", reqData.Name)
	u.RawQuery = q.Encode()

	_, err := http.Get(u.String())
	return err
}

func getNameOfPoint(p int) string {
	if p == 1 {
		return "INN"
	} else if p == 2 {
		return "KZN"
	} else if p == 3 {
		return "AIRPORT_KZN"
	}
	return "BRUH"
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

	trip, err := h.services.Trip.GetById(uctx.UserId, tripId)
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

	newAdminId, err := h.services.Trip.Delete(uctx.UserId, tripId)
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

// func (h *Handler) getCoTravellers(c *gin.Context) {
// 	uctx, err := getUserCtx(c)
// 	if err != nil {
// 		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	tripId, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// }

// func (h *Handler) updateTrip(c *gin.Context) {

// }
