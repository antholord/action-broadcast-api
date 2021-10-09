package websockets

import (
	"encoding/json"
	"fmt"
)

// Session maintains the set of active clients and broadcasts messages to the
// clients.
type Session struct {

	manager *Manager

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

func NewSession(m *Manager, id string) *Session {
	return &Session{
		manager: m,
		sessionId: id,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (s *Session) sendSessionInfo(c *Client, event string) {
	var m *sessionInfoMessage
	switch event {
		case UserJoinEvent:
			m = NewClientJoinedMessage(s, c)
		case UserLeaveEvent:
			m = NewClientLeftMessage(s, c)
		default:
			fmt.Println("Cannot find event message type to create")
			return
	}
	
	bytes, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	s.broadcast <- bytes
}

func (s *Session) getClientNames() []string{
	names := []string{}
	for c := range s.clients {
		names = append(names, c.Name)
	}
	return names
}

func (s *Session) deleteClient(c *Client) {
	close(c.send)
	delete(s.clients, c)
	if len(s.clients) == 0 {
		delete(s.manager.sessions, s.sessionId)
	}
}

func (s *Session) Run() {
	for {
		select {
		case client := <-s.register:
			fmt.Println("Registering client")
			s.clients[client] = true
			go s.sendSessionInfo(client, UserJoinEvent)

		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				fmt.Println("Unregistering client")
				s.deleteClient(client)
				go s.sendSessionInfo(client, UserJoinEvent)
			}

		case message := <-s.broadcast:
			parsedMsg, err := ParseMessage(message)
			fmt.Println(parsedMsg)
			if err != nil {
				return
			}
			for client := range s.clients {
				if parsedMsg.Target == "" || client.Name == parsedMsg.Target {
					fmt.Println("Sending message to 1 client")
					select {
					case client.send <- message:
					default:
						s.deleteClient(client)
					}
				}
			}
		}
	}
}