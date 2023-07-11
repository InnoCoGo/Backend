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
		jt := user.Group("/join_trip")
		{
			jt.POST("/redirect", h.redirectReqToBot)
			jt.PUT("/get_req_from_bot", h.getReqFromBot)
		}
	}
}

func (h *Handler) redirectReqToBot(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

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
	if trip.AdminId != input.AdminId {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	input.UserId = uctx.UserId
	input.TripName = getTripName(trip.FromPoint, trip.ToPoint, trip.ChosenTimestamp)
	input.SecretToken = os.Getenv("BACKEND_SECRET_TOKEN")

	err = doRequest(http.MethodPost, os.Getenv("TG_BOT_URL"), path.Join("/", "join_request"), input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (h *Handler) getReqFromBot(c *gin.Context) { // TODO: webhook
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
	AdminId     int    `json:"trip_admin_id" binding:"required"`
	TripId      int    `json:"trip_id" binding:"required"`
	UserId      int    `json:"id_of_person_asking_to_join"`
	SecretToken string `json:"secret_token"`
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
