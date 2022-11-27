package websockets

import (
	"encoding/json"
)

const (
	UserLeaveEvent    = "user-left"
	GetUsersEvent     = "get-users"
	UsersList = "users-list"
)

type sessionInfoMessage struct {
	SessionId string `json:"sessionId"`
	Type     string `json:"type"`
	Payload   []byte `json:"payload"`
}

func NewClientLeftMessage(session *Session, c *Client) *sessionInfoMessage {
	payload := struct { User string}{ User: c.Name}
	bytes, err := json.Marshal(payload)
	if err == nil {
		println("NewClientLeftMessage error")
		println(err)
	}
	return &sessionInfoMessage{SessionId: session.sessionId, Type: UserLeaveEvent, Payload: bytes}
}

func NewSendUsersListMessage(session *Session) *sessionInfoMessage {
	clients := make([]*Client, 0, len(session.clients))
	for _, value := range clients {
		clients = append(clients, value)
	}
	payload := struct { Users []*Client}{ Users: clients}
	bytes, err := json.Marshal(payload)
	if err == nil {
		println("NewGetUsersMessage error")
		println(err)
	}
	return &sessionInfoMessage{SessionId: session.sessionId, Type: UsersList, Payload: bytes}
}
