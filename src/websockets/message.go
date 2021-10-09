package websockets

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	UserJoinEvent = "user-joined"
	UserLeaveEvent = "user-left"
	ActionEvent = "action"
)

type message struct {
	Target string
	Event string
}

type sessionInfoMessage struct {
	message
	clientNames []string
	userName string
}

func NewMessage() *message {
	return &message{ Target: "", Event: ActionEvent}
}

func NewClientJoinedMessage(s *Session, c *Client) *sessionInfoMessage {
	m := newSessionInfoMessage(s.getClientNames())
	m.Event = UserJoinEvent
	m.userName = c.Name
	return m
}

func NewClientLeftMessage(s *Session, c *Client) *sessionInfoMessage {
	m := newSessionInfoMessage(s.getClientNames())
	m.Event = UserLeaveEvent
	m.userName = c.Name
	return m
}

func newSessionInfoMessage(clientNames []string) *sessionInfoMessage {
	return &sessionInfoMessage { message: message{ Target: ""}, clientNames: clientNames}
}

func ParseMessage(data []byte) (*message, error) {
	msg := &message{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Event == "" {
		fmt.Println("Error reading data, format invalid")
		return nil, errors.New("error reading data, format invalid")
	}
	return msg, nil
}