package p2p

import (
	"errors"
	"strings"
)

var ErrInvalidMessage = errors.New("invalid message")

type Message struct {
	Type string
	Data string
}

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
