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

	processMessage chan *ClientMessage

	// Inbound messages from the clients.
	clientBroadcast chan *message

	serverBroadcast chan *sessionInfoMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewSession(m *Manager, id string) *Session {
	return &Session{
		manager: m,
		sessionId: id,
		processMessage: make(chan *ClientMessage),
		clientBroadcast:  make(chan *message),
		serverBroadcast:  make(chan *sessionInfoMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
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

		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				fmt.Println("Unregistering client")
				s.deleteClient(client)
				s.serverBroadcast <- NewSendUsersListMessage(s)
			}

		case clientMessage := <-s.processMessage:
			parsedMsgData, err := ParseMessage(clientMessage.message)
			fmt.Println(parsedMsgData)
			if err != nil {
				fmt.Println(err)
				return
			}
			switch parsedMsgData.Event {
			case SetDataEvent:
				//TODO see if this works
				clientMessage.client.Data = parsedMsgData.Payload
				s.serverBroadcast <- NewSendUsersListMessage(s)
			case ActionEvent:
				s.clientBroadcast <- parsedMsgData
			}


		case message := <-s.clientBroadcast:
			bytes, _ := json.Marshal(message)
			for client := range s.clients {
				if message.Target == "" || client.Name == message.Target {
					select {
					case client.send <- bytes:
					default:
						s.deleteClient(client)
					}
				}
			}
		case message := <-s.serverBroadcast:
			bytes, _ := json.Marshal(message)
			for client := range s.clients {
				select {
				case client.send <- bytes:
				default:
					s.deleteClient(client)
				}
			}
		}
	}
}