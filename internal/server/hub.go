package server

import (
	"fmt"
	"strconv"

	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/internal/service"
	"github.com/sirupsen/logrus"
)

type Hub struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *core.Message
	Clients    map[int64]*Client `json:"clients"`
}

func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *core.Message, 5),
		Clients:    make(map[int64]*Client),
	}
}

func (hub *Hub) ConsumerKafkaMsg(msg *core.Message) {
	hub.Broadcast <- msg
}

func (h *Hub) Run(s *service.Service) {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Clients[cl.Id]; !ok {
				h.Clients[cl.Id] = cl
			}

		case cl := <-h.Unregister:
			if _, ok := h.Clients[cl.Id]; ok {
				if len(h.Clients) != 0 {
					h.Broadcast <- &core.Message{
						Content:      "User " + cl.Username + " has left the room" + strconv.Itoa(int(cl.RoomId)),
						ContentType:  core.TEXT,
						FromUsername: "System",
						FromUserId:   0,
						ToRoomId:     cl.RoomId,
					}
				}
				delete(h.Clients, cl.Id)
				close(cl.Message)
			}
		case msg := <-h.Broadcast:
			members, err := s.Trip.GetJoinedUsers(msg.ToRoomId)
			if err != nil {
				logrus.Errorf("getting from BROADCAST, GetJoinedUsers() %s", err.Error())
				continue
			}
			fmt.Printf("\nKirdi BROADCAS-qa: %v ---- %v \n", members, msg)

			for _, member := range members {
				if member.UserId == msg.FromUserId || msg.FromUserId == 0 {
					continue
				}
				cl, ok := h.Clients[member.UserId]
				if !ok {
					s.Message.Save(*msg)
					continue
				}

				cl.Message <- msg
			}
			// if _, ok := h.Rooms[m.ToRoomId]; ok {
			// 	_, exits := h.Rooms[m.FromUserId]
			// 	if exits && m.FromUserId != 0 {
			// 		s.Message.Save(*m)
			// 	}
			// 	for _, cl := range h.Rooms[m.ToRoomId].Clients {
			// 		fmt.Printf("\n%v <- %v\n", cl, m)
			// 		cl.Message <- m
			// 	}
			// }
		}
	}
}
