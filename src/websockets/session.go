package websockets

import "fmt"

// Session maintains the set of active clients and broadcasts messages to the
// clients.
type Session struct {

	sessionId string

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewSession(id string) *Session {
	return &Session{
		sessionId: id,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Session) Run() {
	for {
		select {
		case client := <-h.register:
			fmt.Println("Registering client")
			h.clients[client] = true
		case client := <-h.unregister:
			fmt.Println("Unregistering client")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			fmt.Println("Broadcasting a message")
			for client := range h.clients {
				fmt.Println("Sending message to 1 client")
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}