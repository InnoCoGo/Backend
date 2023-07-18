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

	langId := identifyLang(trip.Description)
	// translater := tr.New(tr.Config{
	// 	Url: os.Getenv("TRANSLATE_URL"),
	// 	Key: os.Getenv("TRANSLATE_API_KEY"),
	// })
	// translated, err := translater.Translate(trip.Description, getLang(langId), getLang(!langId))
	// if err != nil {
	// 	response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	translated := ""

	trip.AdminId = uctx.UserId
	trip.AdminUsername = uctx.Username
	trip.AdminTgId = uctx.TgId

	trip.TranslatedDesc = translated
	if langId {
		trip.Description, trip.TranslatedDesc = trip.TranslatedDesc, trip.Description
	}

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

func getTripName(from, to int, timestamp string) string {
	return fmt.Sprintf("%d -> %d at: %s", from, to, timestamp)
}

func identifyLang(text string) bool {
	for _, r := range text {
		if r > 127 {
			return false
		} else if 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' {
			return true
		}
	}
	return true
}

func getLang(langId bool) string {
	if langId {
		return "en"
	}
	return "ru"
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

	tripId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trip, err := h.services.Trip.GetById(tripId)
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

	tripId, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

func (h *Handler) getJoinedUsers(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tripId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.services.Trip.GetJoinedUsers(uctx.UserId, tripId)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: users})
}
