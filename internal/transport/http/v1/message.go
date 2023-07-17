package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/response"
)

func (h *Handler) initMessagesRoutes(api *gin.RouterGroup) {
	messages := api.Group("/messages", h.userIdentity)
	{
		messages.GET("/:room_id", h.fetchRoomMessages)
	}
}

func (h *Handler) fetchRoomMessages(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Param("room_id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	messages, err := h.services.Message.FetchRoomMessages(int64(roomId))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: messages})
}
