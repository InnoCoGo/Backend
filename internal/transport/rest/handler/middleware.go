package handler

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var BotToken = os.Getenv("BOT_TOKEN")

const (
	AuthorizationHeader = "Authorization"
	userCtx             = "userId"
	webappKeyword       = "WebAppData"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(AuthorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty header!")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "incorrect passing of header!")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

// func (h *Handler) tgUserIdentity(c *gin.Context) {

// 	url := c.Request.URL
// 	token := url.Query().Get("token")
// 	if len(token) == 0 {
// 		newErrorResponse(c, http.StatusUnauthorized, "empty tg token")
// 		return
// 	}

// 	userId, err := h.services.Authorization.ParseTgToken(token) // TODO!!!
// 	if err != nil {
// 		newErrorResponse(c, http.StatusUnauthorized, err.Error())
// 		return
// 	}

// 	c.Set(userCtx, userId)
// }

// func getUserIdentity(c *gin.Context) (int, string, error) {
// 	userId, ok := c.Get(userCtx)
// }

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)

	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, errors.New("user id not found")
	}

	return idInt, nil
}
