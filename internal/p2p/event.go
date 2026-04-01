package p2p

import (
	"encoding/json"
	"log"
)

type PeerEventType string

type PeerEvent struct {
	Type     PeerEventType `json:"type"`
	PeerLink string        `json:"link"`
	Message  string        `json:"msg"`
}

// Create a new peer event
func NewPeerEvent(eventType PeerEventType, peerLink string, msg string) *PeerEvent {
	return &PeerEvent{
		Type:     eventType,
		PeerLink: peerLink,
		Message:  msg,
	}
}

// Convert a peer event to JSON format
func (pe *PeerEvent) Json() string {
	j, err := json.Marshal(pe)
	if err != nil {
		log.Panic(err)
	}
	return string(j)
}

const EventFeedLimit = 10

type EventFeed struct {
	subscribers []chan *PeerEvent
}

func (ef *EventFeed) Subscribe() chan *PeerEvent {
	ch := make(chan *PeerEvent, EventFeedLimit)
	ef.subscribers = append(ef.subscribers, ch)
	return ch
}

func (ef *EventFeed) Send(event *PeerEvent) chan *PeerEvent {
	for _, ch := range ef.subscribers {
		select {
		case ch <- event:
		default:
			log.Println("Skipped event for slow subscriber")
		}
	}

	return nil
}
