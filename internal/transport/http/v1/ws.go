package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/itoqsky/InnoCoTravel-backend/internal/server"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/response"
)

func (h *Handler) initWsRoutes(api *gin.RouterGroup) {
	ws := api.Group("/ws")
	{
		ws.GET("/join_trip/:trip_id", h.joinTrip)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) joinTrip(c *gin.Context) {
	uctx, err := getUserCtx(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	tripId, err := strconv.Atoi(c.Param("trip_id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cl := &server.Client{
		Conn:     conn,
		Message:  make(chan *server.Message, 10),
		Id:       uctx.UserId,
		RoomId:   int64(tripId),
		Username: uctx.Username,
	}

	m := &server.Message{
		Content:  fmt.Sprintf(`%s has joined the trip`, uctx.Username),
		RoomId:   int64(tripId),
		Username: uctx.Username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.WriteMessage()
	cl.ReadMessage(h.hub)
}
