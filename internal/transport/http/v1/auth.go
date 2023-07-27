package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/itoqsky/InnoCoTravel-backend/docs"
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/response"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)

		auth.POST("/tg-login", h.tgLogIn)
	}
}

func (h *Handler) signUp(c *gin.Context) {
	var user core.User
	if err := c.BindJSON(&user); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (h *Handler) signIn(c *gin.Context) {
	var userSignInObj signInInput

	if err := c.BindJSON(&userSignInObj); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := core.User{
		Username:       userSignInObj.Username,
		PasswordOrHash: userSignInObj.Password,
	}

	userId, err := h.services.Authorization.GetUserId(user)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(core.UserCtx{UserId: userId, Username: user.Username})
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{token})
}

type TgLoginInput struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoUrl  string `json:"photo_url"`
	QueryId   string `json:"query_id"`
	User      string `json:"user"`
	AuthDate  int    `json:"auth_date" binding:"required"`
	Hash      string `json:"hash" binding:"required"`
}

type tgLoginWApp struct {
	User     string `json:"user" binding:"required"`
	AuthDate int    `json:"auth_date" binding:"required"`
	QueryId  string `json:"query_id" binding:"required"`
	Hash     string `json:"hash" binding:"required"`
}

type tgLoginWSite struct {
	Id        int64  `json:"id" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Username  string `json:"username" binding:"required"`
	PhotoUrl  string `json:"photo_url" binding:"required"`
	AuthDate  int    `json:"auth_date" binding:"required"`
	Hash      string `json:"hash" binding:"required"`
}

const (
	webAppKeyword = "WebAppData"
)

func (h *Handler) tgLogIn(c *gin.Context) {
	var (
		rawTgUser TgLoginInput
		tgUserWA  tgLoginWApp
		tgUserWS  tgLoginWSite
		user      core.User
	)

	var jsonData []byte
	var authData map[string]interface{}

	var keyword string
	var err error

	if err = c.BindJSON(&rawTgUser); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if rawTgUser.User == "" {
		user = core.User{
			FirstName: rawTgUser.FirstName,
			LastName:  rawTgUser.LastName,
			Username:  rawTgUser.Username,
			TgId:      rawTgUser.Id,
		}

		jsonData, _ = json.Marshal(rawTgUser)
		err = json.Unmarshal(jsonData, &tgUserWS)
		if err != nil {
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		jsonData, _ = json.Marshal(tgUserWS)
		err = json.Unmarshal(jsonData, &authData)
	} else {
		var userField map[string]interface{}
		err = json.Unmarshal([]byte(rawTgUser.User), &userField)
		if err != nil {
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		for k, v := range userField {
			if k != "first_name" && k != "last_name" && k != "username" && k != "id" && k != "language_code" && k != "allows_write_to_pm" {
				response.NewErrorResponse(c, http.StatusBadRequest, "incorrect keys of user field from telegram webapp")
				return
			}
			if k == "id" {
				if _, ok := v.(float64); !ok {
					response.NewErrorResponse(c, http.StatusBadRequest, "incorrect value of id in user field from telegram webapp")
					return
				}
				userField[k] = int64(v.(float64))
			} else {
				if _, ok := v.(string); !ok {
					response.NewErrorResponse(c, http.StatusBadRequest, "incorrect assertion string in user field from telegram webapp")
					return
				}
			}
		}
		user = core.User{
			FirstName: userField["first_name"].(string),
			LastName:  userField["last_name"].(string),
			Username:  userField["username"].(string),
			TgId:      userField["id"].(int64),
		}
		keyword = webAppKeyword

		jsonData, _ = json.Marshal(rawTgUser)
		err = json.Unmarshal(jsonData, &tgUserWA)
		if err != nil {
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		jsonData, _ = json.Marshal(tgUserWA)
		err = json.Unmarshal(jsonData, &authData)
	}
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ok, err := h.services.Authorization.VerifyTgAuthData(authData, keyword)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !ok {
		response.NewErrorResponse(c, http.StatusBadRequest, "Data is NOT from telegram")
		return
	}

	userId, err := h.services.Authorization.GetUserId(user)
	if err != nil {
		userId, err = h.services.Authorization.CreateUser(user)
		if err != nil {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	token, err := h.services.Authorization.GenerateToken(core.UserCtx{UserId: userId, Username: user.Username, TgId: user.TgId})
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{token})
}
