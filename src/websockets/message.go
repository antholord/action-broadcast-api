package websockets

import (
	"encoding/json"
	"errors"
	"fmt"
)

type message struct {
	Target  string 	`json:"target"`
	Event   string 	`json:"event"`
	Payload map[string]interface{} `json:"payload"`
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