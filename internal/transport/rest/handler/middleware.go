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
)

var BotToken = os.Getenv("BOT_TOKEN")

const (
	AuthorizationHeader = "Authorization"
	userCtx             = "userId"
)

// func (h *Handler) userIdentity(c *gin.Context) {
// 	header := c.GetHeader(AuthorizationHeader)
// 	if header != "" {
// 		newErrorResponse(c, http.StatusUnauthorized, "empty header!")
// 		return
// 	}
// 	headerParts := strings.Split(header, " ")
// 	if len(headerParts) != 2 {
// 		newErrorResponse(c, http.StatusUnauthorized, "incorrect passing of header!")
// 		return
// 	}

// 	userId, err := h.services.Authorization.ParseToken(headerParts[1])
// 	if err != nil {
// 		newErrorResponse(c, http.StatusUnauthorized, err.Error())
// 		return
// 	}

// 	c.Set(userCtx, userId)
// }

type TgUserWebSite struct {
	Id        int     `json:"id" binding:"required"`
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Username  string  `json:"username" binding:"required"`
	PhotoUrl  string  `json:"photo_url" binding:"required"`
	AuthDate  float64 `json:"auth_date" binding:"required"`
	Hash      string  `json:"hash" binding:"required"`
}

type TgUserWebApp struct {
	User     string  `json:"user"`
	AuthDate float64 `json:"auth_date"`
	QueryId  string  `json:"query_id"`
	Hash     string  `json:"hash"`
}

func (h *Handler) tgUserAuthorize(c *gin.Context) {
	authData, secretKey, userId, err := parseInput(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	checkHash, _ := authData["hash"].(string)

	delete(authData, "hash")

	var dataCheckArr []string
	for key, val := range authData {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%v", key, val))
	}
	sort.Strings(dataCheckArr)

	dataCheckString := strings.Join(dataCheckArr, "\n")

	hmc := hmac.New(sha256.New, secretKey[:])
	hmc.Write([]byte(dataCheckString))
	hash := hmc.Sum(nil)

	if !hmac.Equal(hash, []byte(checkHash)) {
		newErrorResponse(c, http.StatusUnauthorized, "Data is NOT from Telegram")
		return
	}

	authDate, ok := authData["auth_date"].(float64)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid authorization data")
		return
	}
	if float64(time.Now().Unix()-int64(authDate)) > 86400 {
		newErrorResponse(c, http.StatusUnauthorized, "Data is outdated")
		return
	}

	// 															TODO ! (Chech if user in database)

	c.Set(userCtx, userId)
}

func parseInput(c *gin.Context) (map[string]interface{}, []byte, int, error) {
	var tguserWA TgUserWebApp
	var tguserWS TgUserWebSite

	authData := make(map[string]interface{})

	secretKey := make([]byte, 32)

	userId := 0

	if err := c.BindJSON(&tguserWA); err != nil {
		if err = c.BindJSON(&tguserWS); err != nil {
			return nil, secretKey, 0, err
		}

		jsonData, _ := json.Marshal(tguserWS)
		json.Unmarshal(jsonData, &authData)

		h := hmac.New(sha256.New, []byte("WebAppData"))
		h.Write([]byte(BotToken))
		secretKey := h.Sum(nil)

		userId = authData["id"].(int)
		return authData, secretKey, userId, nil
	}

	jsonData, _ := json.Marshal(tguserWA)
	json.Unmarshal(jsonData, &authData)

	// secretKey = sha256.Sum256([]byte(BotToken))				TODO!
	h := hmac.New(sha256.New, nil)
	h.Write([]byte(BotToken))
	secretKey = h.Sum(nil)

	// var user core.User 									 	TODO!
	var str = authData["user"].(string)
	var err error

	if len(str) != 0 {
		for l := 0; l < len(str); l++ {
			r := l
			for r != ':' {
				r++
			}
			r++
			l = r
			for r != ',' {
				userId *= 10
				userId += int(str[l])
			}
		}
	} else {
		err = fmt.Errorf("Error, empty user key in the reqest!")
	}

	return authData, secretKey, userId, err
}
