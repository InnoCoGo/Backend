package server

type Room struct {
	Id      int64             `json:"room_id"`
	Name    string            `json:"room_name"`
	Clients map[int64]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[int64]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[int64]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (hub *Hub) ConsumerKafkaMsg(msg *Message) {
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
						h.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomId:   cl.RoomId,
							Username: cl.Username,
						}
					}
					delete(h.Rooms[cl.RoomId].Clients, cl.Id)
					close(cl.Message)
				}
			}
		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.RoomId]; ok {
				for _, cl := range h.Rooms[m.RoomId].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
