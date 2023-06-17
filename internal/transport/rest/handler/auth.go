package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
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

type tgUserWebApp struct {
	User     string `json:"user" binding:"required"`
	AuthDate int    `json:"auth_date" binding:"required"`
	QueryId  string `json:"query_id" binding:"required"`
	Hash     string `json:"hash" binding:"required"`
}

//	type userFieldWA struct {
//		Id           int    `json:"id" binding:"required"`
//		FirstName    string `json:"first_name" binding:"required"`
//		LastName     string `json:"last_name" binding:"required"`
//		Username     string `json:"username" binding:"required"`
//		LanguageCode string `json:"language_code" binding:"required"`
//	}
type tgUserWebSite struct {
	Id        int    `json:"id" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Username  string `json:"username" binding:"required"`
	PhotoUrl  string `json:"photo_url" binding:"required"`
	AuthDate  int    `json:"auth_date" binding:"required"`
	Hash      string `json:"hash" binding:"required"`
}

type TgUserCredentials struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoUrl  string `json:"photo_url"`
	QueryId   string `json:"query_id"`
	User      string `json:"user"`
	AuthDate  int    `json:"auth_date" binding:"required"`
	Hash      string `json:"hash" binding:"required"`
}

const (
	webAppKeyword = "WebAppData"
)

func (h *Handler) tgLogIn(c *gin.Context) {
	var (
		rawTgUser TgUserCredentials
		tgUserWA  tgUserWebApp
		tgUserWS  tgUserWebSite
		user      core.User
	)

	var jsonData []byte
	var authData map[string]interface{}

	var keyword string
	var err error

	if err = c.BindJSON(&rawTgUser); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
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
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		jsonData, _ = json.Marshal(tgUserWS)
		err = json.Unmarshal(jsonData, &authData)
	} else {
		var userField map[string]interface{}
		err = json.Unmarshal([]byte(rawTgUser.User), &userField)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		for k, v := range userField {
			if k != "first_name" && k != "last_name" && k != "username" && k != "id" && k != "language_code" {
				newErrorResponse(c, http.StatusBadRequest, "incorrect keys of user field from telegram webapp")
				return
			}
			if k == "id" {
				if _, ok := v.(float64); !ok {
					newErrorResponse(c, http.StatusBadRequest, "incorrect value of id in user field from telegram webapp")
					return
				}
				userField[k] = int(v.(float64))
			} else {
				if _, ok := v.(string); !ok {
					newErrorResponse(c, http.StatusBadRequest, "incorrect assertion string in user field from telegram webapp")
					return
				}
			}
		}
		user = core.User{
			FirstName: userField["first_name"].(string),
			LastName:  userField["last_name"].(string),
			Username:  userField["username"].(string),
			TgId:      userField["id"].(int),
		}
		keyword = webAppKeyword

		jsonData, _ = json.Marshal(rawTgUser)
		err = json.Unmarshal(jsonData, &tgUserWA)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		jsonData, _ = json.Marshal(tgUserWA)
		err = json.Unmarshal(jsonData, &authData)
	}
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

func verifyTgAuthData(authData map[string]interface{}, keyword string) (bool, error) {
	checkHash, _ := authData["hash"].(string)
	authData["auth_date"] = int(authData["auth_date"].(float64))

	authIdWS, ok := authData["id"].(float64)
	if ok {
		authData["id"] = int(authIdWS)
	}

	delete(authData, "hash")

	var dataCheckArr []string
	for key, val := range authData {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%v", key, val)) // WARNING!
	}
	sort.Strings(dataCheckArr)
	dataCheckString := strings.Join(dataCheckArr, "\n")

	var h hash.Hash
	if keyword == "" {
		h = sha256.New()
	} else {
		h = hmac.New(sha256.New, []byte(keyword))
	}
	h.Write([]byte(os.Getenv("BOT_TOKEN")))
	secretKey := h.Sum(nil)

	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	hash := hex.EncodeToString(h.Sum(nil))

	if hash != checkHash {
		return false, fmt.Errorf("the hashes don't match")
	}

	authDate, _ := authData["auth_date"].(int)
	if time.Now().Unix()-int64(authDate) > 86400 {
		return false, fmt.Errorf("expired auth date")
	}

	return true, nil
}
