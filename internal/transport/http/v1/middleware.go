package v1

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var BotToken = os.Getenv("BOT_TOKEN")

const (
	AuthorizationHeader = "Authorization"
	userCtx             = "userId"
	tgCtx               = "tgId"
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

	id, err := h.services.Authorization.ParseToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, id.User)
	// c.Set(tgCtx, id.TG)
}

func getUserId(c *gin.Context) int {
	id, ok := c.Get(userCtx)

	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0
	}

	return idInt
}

func getTgId(c *gin.Context) int { // for grpc client communication with telegram bot
	id, ok := c.Get(tgCtx)

	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "tg id not found")
		return 0
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "tg id is of invalid type")
		return 0
	}

	return idInt
}
