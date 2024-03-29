package server

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/itoqsky/InnoCoTravel-backend/internal/kafka"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/protocol"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *protocol.Message
	Id       int64  `json:"client_id"`
	Username string `json:"username"`
	RoomId   int64  `json:"room_id"`
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

		msg := &protocol.Message{
			Content:      string(msgPack),
			FromUsername: c.Username,
			FromId:       c.Id,
			ToRoomId:     c.RoomId,
		}
		log.Print(msg)

		if viper.GetBool("kafka.enabled") {
			kafka.Produce(msgPack)
		} else {
			hub.Broadcast <- msg
		}
	}
}
