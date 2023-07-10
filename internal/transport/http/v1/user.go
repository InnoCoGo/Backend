package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/response"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	user := api.Group("/user")
	{
		user.POST("/sendReqJoinTrip", h.joinTripSendReq)
		user.GET("/get_req_from_ bot", h.getReqFromBot)
	}
}

func (h *Handler) joinTripSendReq(c *gin.Context) {
	// uctx, err := getUserCtx(c)
	// if err != nil {
	// 	response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	var input joinRequest

	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trip, err := h.services.Trip.GetById(input.TripId)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	input.TripName = getTripName(trip.FromPoint, trip.ToPoint, trip.ChosenTimestamp)

	err = doRequest(path.Join("/", "join_request"), os.Getenv("TG_BOT_URL"), input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) getReqFromBot(c *gin.Context) {
	var input joinRequest
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.SecretToken != os.Getenv("BACKEND_SECRET_TOKEN") {
		response.NewErrorResponse(c, http.StatusBadRequest, "wrong secret token!")
		return
	}

	if err := h.services.User.JoinTrip(input.UserId, input.TripId); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

type joinRequest struct {
	AdminId     int    `json:"trip_admin_id" binding:"required"`
	TripId      int    `json:"trip_id" binding:"required"`
	UserId      int    `json:"id_of_person_asking_to_join" binding:"required"`
	SecretToken string `json:"secret_token"`
	TripName    string `json:"trip_name"`
}

func doRequest(host, p string, join_req_body joinRequest) error {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   p,
	}

	reqbody, err := json.Marshal(join_req_body)

	if err != nil {
		return err
	}

	res, err := http.Post(u.String(), "application/json", bytes.NewBuffer(reqbody))
	defer func() { _ = res.Body.Close() }()
	return err
}
