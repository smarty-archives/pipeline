package messaging

import "time"

type Dispatch struct {
	SourceID    uint64
	MessageID   uint64
	Destination string
	MessageType string
	Encoding    string
	Durable     bool
	Expiration  time.Time
	Payload     []byte
	Message     interface{}
}
