package v1

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
)

var BotToken = os.Getenv("BOT_TOKEN")

const (
	AuthorizationHeader = "Authorization"
	userIdCtx           = "userId"
	usernameCtx         = "username"
	webappKeyword       = "WebAppData"
)

func (h *Handler) userIdentity(c *gin.Context) {
	var token string
	header := c.GetHeader(AuthorizationHeader)
	if header != "" {
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			newErrorResponse(c, http.StatusUnauthorized, "incorrect passing of header!")
			return
		}
		token = headerParts[1]
	} else {
		url := c.Request.URL
		token = url.Query().Get("token")
		if len(token) == 0 {
			newErrorResponse(c, http.StatusUnauthorized, "empty url token param")
			return
		}
	}

	uctx, err := h.services.Authorization.ParseToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userIdCtx, uctx.UserId)
	c.Set(usernameCtx, uctx.Username)
}

func getUserCtx(c *gin.Context) (core.UserCtx, error) {
	id, ok := c.Get(userIdCtx)

	if !ok {
		return core.UserCtx{}, fmt.Errorf("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return core.UserCtx{}, fmt.Errorf("user id is of invalid type")
	}

	username, ok := c.Get(usernameCtx)

	if !ok {
		return core.UserCtx{}, fmt.Errorf("(tg)username not found")
	}
	usernameStr, ok := username.(string)
	if !ok {
		return core.UserCtx{}, fmt.Errorf("(tg)username is of invalid type")
	}

	log.Printf("\n%v && %v\n", idInt, usernameStr)

	return core.UserCtx{UserId: idInt, Username: usernameStr}, nil
}
