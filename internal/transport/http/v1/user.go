package v1

import (
	"bytes"
	"crypto/tls"
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
	user := api.Group("/user")
	{
		jt := user.Group("/join_trip")
		{
			jt.GET("/req/:trip_id", h.userIdentity, h.redirectReqToBot)
			jt.POST("/res", h.getResFromBot)
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

	trip, err := h.services.Trip.GetById(int64(trip_id))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var redirectReq = InputGetResFromBot{
		AdminTgId:   trip.AdminTgId,
		UserId:      uctx.UserId,
		UserTgId:    uctx.TgId,
		TripId:      int64(trip_id),
		SecretToken: os.Getenv("BACKEND_SECRET_TOKEN"),
		TripName:    getTripName(trip.FromPoint, trip.ToPoint, trip.ChosenTimestamp),
	}

	err = doRequest(http.MethodPost, os.Getenv("TG_BOT_URL"), path.Join("/", "join_request"), redirectReq, true)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (h *Handler) getResFromBot(c *gin.Context) { // TODO: webhook
	var input InputGetResFromBot
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

type InputGetResFromBot struct {
	TripId      int64  `json:"trip_id" binding:"required"`
	AdminTgId   int64  `json:"trip_admin_tg_id"`
	UserId      int64  `json:"id_of_person_asking_to_join" binding:"required"`
	UserTgId    int64  `json:"tg_id_of_person_asking_to_join"`
	SecretToken string `json:"secret_token" binding:"required"`
	Accepted    bool   `json:"accepted" binding:"required"`
	TripName    string `json:"trip_name"`
}

func doRequest(methd, host, p string, bodyStruct interface{}, secure bool) error {
	s := "http"
	if secure {
		s = "https"
	}
	u := url.URL{
		Scheme: s,
		Host:   host,
		Path:   p,
	}

	body, err := json.Marshal(bodyStruct)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(methd, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	httpCl := http.Client{}

	if secure {
		httpCl = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	res, err := httpCl.Do(req)
	defer func() {
		if res != nil {
			_ = res.Body.Close()
		}
	}()
	return err
}
