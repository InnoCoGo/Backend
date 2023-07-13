package v1

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/response"
)

var BotToken = os.Getenv("BOT_TOKEN")

const (
	AuthorizationHeader = "Authorization"
	userIdCtx           = "userId"
	userTgIdCtx         = "userId"
	usernameCtx         = "username"
	webappKeyword       = "WebAppData"
)

func (h *Handler) userIdentity(c *gin.Context) {
	var token string
	header := c.GetHeader(AuthorizationHeader)
	if header != "" {
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			response.NewErrorResponse(c, http.StatusUnauthorized, "incorrect passing of header!")
			return
		}
		token = headerParts[1]
	} else {
		url := c.Request.URL
		token = url.Query().Get("token")
		if len(token) == 0 {
			response.NewErrorResponse(c, http.StatusUnauthorized, "empty url token param")
			return
		}
	}

	uctx, err := h.services.Authorization.ParseToken(token)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userIdCtx, uctx.UserId)
	c.Set(userTgIdCtx, uctx.TgId)
	c.Set(usernameCtx, uctx.Username)
}

func getUserCtx(c *gin.Context) (core.UserCtx, error) {
	id, ok := c.Get(userIdCtx)
	if !ok {
		return core.UserCtx{}, fmt.Errorf("user id not found")
	}
	idInt, ok := id.(int64)
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

	tgId, ok := c.Get(userTgIdCtx)
	if !ok {
		return core.UserCtx{}, fmt.Errorf("user id not found")
	}
	tgIdInt, ok := tgId.(int64)
	if !ok {
		return core.UserCtx{}, fmt.Errorf("user id is of invalid type")
	}

	log.Printf("\n%v && %v && %v\n", idInt, usernameStr, tgIdInt)

	return core.UserCtx{UserId: idInt, Username: usernameStr, TgId: tgIdInt}, nil
}
