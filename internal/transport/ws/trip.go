package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/itoqsky/InnoCoTravel-backend/internal/server"
)

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) createRoom(c *gin.Context) {
	req_id := c.Query("id")
	req_name := c.Query("name")

	h.hub.Rooms[req_id] = &server.Room{
		Id:      req_id,
		Name:    req_name,
		Clients: make(map[string]*server.Client),
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":   req_id,
		"name": req_name,
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) joinRoom(c *gin.Context) {
	roomID := c.Param("tripId")
	clientID := c.Query("userId")
	username := c.Query("username")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cl := &server.Client{
		Conn:     conn,
		Message:  make(chan *server.Message, 10),
		Id:       clientID,
		RoomId:   roomID,
		Username: username,
	}

	m := &server.Message{
		Content:  "A new user has joined the room",
		RoomId:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.WriteMessage()
	cl.ReadMessage(h.hub)
}
