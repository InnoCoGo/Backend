package server

import "github.com/itoqsky/InnoCoTravel-backend/pkg/protocol"

type Room struct {
	Id      int64             `json:"room_id"`
	Name    string            `json:"room_name"`
	Clients map[int64]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[int64]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *protocol.Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[int64]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *protocol.Message, 5),
	}
}

func (hub *Hub) ConsumerKafkaMsg(msg *protocol.Message) {
	hub.Broadcast <- msg
}

func (h *Hub) Run() {
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
						h.Broadcast <- &protocol.Message{
							Content:      "User " + cl.Username + " has left the chat",
							FromUsername: "System",
							FromId:       0,
							ToRoomId:     cl.RoomId,
						}
					}
					delete(h.Rooms[cl.RoomId].Clients, cl.Id)
					close(cl.Message)
				}
			}
		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.ToRoomId]; ok {
				for _, cl := range h.Rooms[m.ToRoomId].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
