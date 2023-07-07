package server

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	Id       string `json:"client_id"`
	Username string `json:"username"`
	RoomId   string `json:"room_id"`
}

type Message struct {
	Content  string `json:"content"`
	RoomId   string `json:"room_id"`
	Username string `json:"username"`
}

func (c *Client) WriteMessage() {
	defer func() { c.Conn.Close() }()
	for {
		msg, ok := <-c.Message
		if !ok {
			return
		}
		c.Conn.WriteJSON(msg)
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msgPack, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Error(err.Error())
			}
			break
		}

		msg := &Message{
			Content:  string(msgPack),
			RoomId:   c.RoomId,
			Username: c.Username,
		}
		hub.Broadcast <- msg
	}
}
