package p2p

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
)

var ErrInvalidMessage = errors.New("invalid message")

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// ParseMessage parses a message based on the message Type
// and the Data contained in the message.
func ParseMessage(line string) (*Message, error) {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return nil, ErrInvalidMessage
	}

	msg := &Message{
		Type: parts[0],
		Data: parts[1],
	}

	return msg, nil
}

func (m *Message) Json() string {
	j, err := json.Marshal(m)
	if err != nil {
		log.Panic(err)
	}
	return string(j)
}
