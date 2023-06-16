package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel_backend/internal/core"
)

func (h *Handler) signUp(c *gin.Context) {
	var user core.User
	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type userSignIn struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var userSignInObj userSignIn

	if err := c.BindJSON(&userSignInObj); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := core.User{
		Username:       userSignInObj.Username,
		PasswordOrHash: userSignInObj.Password,
	}

	userId, err := h.services.Authorization.GetUserId(user)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	token, err := h.services.Authorization.GenerateToken(userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

type userTgWebApp struct {
	User     userFieldWA `json:"user" binding:"required"`
	AuthDate float64     `json:"auth_date" binding:"required"`
	QueryId  string      `json:"query_id" binding:"required"`
	Hash     string      `json:"hash" binding:"required"`
}

type userFieldWA struct {
	Id           int    `json:"id" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Username     string `json:"username" binding:"required"`
	LanguageCode string `json:"language_code" binding:"required"`
}
type userTgWebSite struct {
	Id        int     `json:"id" binding:"required"`
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Username  string  `json:"username" binding:"required"`
	PhotoUrl  string  `json:"photo_url" binding:"required"`
	AuthDate  float64 `json:"auth_date" binding:"required"`
	Hash      string  `json:"hash" binding:"required"`
}

const (
	webAppKeyword = "WebAppData"
)

func (h *Handler) tgLogIn(c *gin.Context) {
	var userWA userTgWebApp
	var userWS userTgWebSite

	var user core.User
	var jsonData []byte
	authData := make(map[string]interface{})

	var keyword []byte

	if err := c.BindJSON(&userWA); err != nil {
		if err = c.BindJSON(&userWS); err != nil {
			newErrorResponse(c, http.StatusBadRequest, "couldn't bind to user (webapp or website)")
			return
		}
		user = core.User{
			FirstName: userWS.FirstName,
			LastName:  userWS.LastName,
			Username:  userWS.Username,
			TgId:      userWS.Id,
		}
		jsonData, _ = json.Marshal(userWS)

	} else {
		user = core.User{
			FirstName: userWA.User.FirstName,
			LastName:  userWA.User.LastName,
			Username:  userWA.User.Username,
			TgId:      userWA.User.Id,
		}
		jsonData, _ = json.Marshal(userWA)
		keyword = []byte(webAppKeyword)
	}
	err := json.Unmarshal(jsonData, &authData)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ok, err := verifyTgAuthData(authData, keyword)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "Data is NOT from telegram")
		return
	}

	id, err := h.services.Authorization.GetUserId(user)
	if err != nil {
		id, err = h.services.Authorization.CreateUser(user)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	token, err := h.services.Authorization.GenerateToken(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func verifyTgAuthData(authData map[string]interface{}, keyword []byte) (bool, error) {
	h := hmac.New(sha256.New, keyword)
	h.Write([]byte(os.Getenv("BOT_TOKEN")))
	secretKey := h.Sum(nil)

	checkHash, _ := authData["hash"].(string)
	delete(authData, "hash")

	var dataCheckArr []string
	for key, val := range authData {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%v", key, val)) // WARNING!
	}
	sort.Strings(dataCheckArr)
	dataCheckString := strings.Join(dataCheckArr, "\n")

	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	hash := h.Sum(nil)

	if !hmac.Equal(hash, []byte(checkHash)) {
		return false, fmt.Errorf("the hashes don't match")
	}

	authDate, ok := authData["auth_date"].(float64)
	if !ok {
		return false, fmt.Errorf("couldn't extract auth date")
	}

	if float64(time.Now().Unix()-int64(authDate)) > 86400 {
		return false, fmt.Errorf("expired auth date")
	}

	return true, nil
}
