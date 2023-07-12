package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/response"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	user := api.Group("/user", h.userIdentity)
	{
		jt := user.Group("/join_trip")
		{
			jt.POST("/req/:trip_id", h.redirectReqToBot)
			jt.PUT("/res", h.getResFromBot)
		}
	}
}

func (h *Handler) redirectReqToBot(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	trip_id, err := strconv.Atoi(c.Param("trip_id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	trip, err := h.services.Trip.GetById(trip_id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var redirectReq = joinRequest{
		UserId:      uctx.UserId,
		AdminId:     trip.AdminId,
		TripId:      trip_id,
		SecretToken: os.Getenv("BACKEND_SECRET_TOKEN"),
		TripName:    getTripName(trip.FromPoint, trip.ToPoint, trip.ChosenTimestamp),
	}

	err = doRequest(http.MethodPost, os.Getenv("TG_BOT_URL"), path.Join("/", "join_request"), redirectReq)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (h *Handler) getResFromBot(c *gin.Context) { // TODO: webhook
	var input joinRequest
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.SecretToken != os.Getenv("BACKEND_SECRET_TOKEN") {
		response.NewErrorResponse(c, http.StatusBadRequest, "wrong secret token!")
		return
	}

	if input.Accepted {
		if err := h.services.User.JoinTrip(input.UserId, input.TripId); err != nil {
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}

type joinRequest struct {
	AdminId     int    `json:"trip_admin_id"`
	TripId      int    `json:"trip_id" binding:"required"`
	UserId      int    `json:"id_of_person_asking_to_join" binding:"required"`
	SecretToken string `json:"secret_token" binding:"required"`
	Accepted    bool   `json:"accepted"`
	TripName    string `json:"trip_name"`
}

func doRequest(methd, host, p string, bodyStruct interface{}) error {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   p,
	}

	body, err := json.Marshal(bodyStruct)
	if err != nil {
		return err
	}

	httpCl := http.Client{}
	req, err := http.NewRequest(methd, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := httpCl.Do(req)

	defer func() {
		if res != nil {
			_ = res.Body.Close()
		}
	}()
	return err
}
