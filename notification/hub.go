package notification

import (
	"First/model"
	"log"
)

type Hub struct {
	Clients    map[int]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan model.Notification
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan model.Notification),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			log.Printf("Registering client with userID: %d\n", client.UserID)
			h.Clients[client.UserID] = client
			h.logConnectedUsers()

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				log.Printf("Unregistering client with userID: %d\n", client.UserID)
				delete(h.Clients, client.UserID)
				close(client.Send)
				h.logConnectedUsers()
			}

		case notification := <-h.Broadcast:
			log.Printf("Broadcasting to user %d", notification.ToUser)
			h.logConnectedUsers()
			if client, ok := h.Clients[notification.ToUser]; ok {
				select {
				case client.Send <- notification:
				default:
					close(client.Send)
					delete(h.Clients, client.UserID)
					h.logConnectedUsers()
				}
			}
		}
	}
}

func (h *Hub) logConnectedUsers() {
	var ids []int
	for id := range h.Clients {
		ids = append(ids, id)
	}
	log.Printf("Connected user IDs: %v\n", ids)
}
