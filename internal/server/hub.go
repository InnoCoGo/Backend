package server

type Room struct {
	Id      int64              `json:"room_id"`
	Name    string             `json:"room_name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}
