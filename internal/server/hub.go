package server

import (
	"fmt"
	"strconv"

	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/internal/service"
)

type Room struct {
	Id      int64             `json:"room_id"`
	Name    string            `json:"room_name"`
	Clients map[int64]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[int64]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *core.Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[int64]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *core.Message, 5),
	}
}

func (hub *Hub) ConsumerKafkaMsg(msg *core.Message) {
	hub.Broadcast <- msg
}

func (h *Hub) Run(s *service.Service) {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Rooms[cl.RoomId]; ok {
				r := h.Rooms[cl.RoomId]
				// if _, ok := r.Clients[cl.Id]; !ok {
				r.Clients[cl.Id] = cl
				// }
			}

		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomId]; ok {
				clients := h.Rooms[cl.RoomId].Clients
				if _, ok := clients[cl.Id]; ok {
					if len(clients) != 0 {
						h.Broadcast <- &core.Message{
							Content:      "User " + cl.Username + " has left the room" + strconv.Itoa(int(cl.RoomId)),
							ContentType:  core.TEXT,
							FromUsername: "System",
							FromUserId:   0,
							ToRoomId:     cl.RoomId,
						}
					}
					delete(h.Rooms[cl.RoomId].Clients, cl.Id)
					close(cl.Message)
				}
			}
		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.ToRoomId]; ok {
				_, exits := h.Rooms[m.FromUserId]
				if exits && m.FromUserId != 0 {
					s.Message.Save(*m)
				}
				for _, cl := range h.Rooms[m.ToRoomId].Clients {
					fmt.Printf("\n%v <- %v\n", cl, m)
					cl.Message <- m
				}
			}
		}
	}
}
