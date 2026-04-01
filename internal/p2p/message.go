package p2p

type Message struct {
	PeerLink string
	Data     string
}

func NewMessage(peerLink string, data string) *Message {
	return &Message{
		PeerLink: peerLink,
		Data:     data,
	}
}
