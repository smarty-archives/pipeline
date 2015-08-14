package messaging

import "time"

type Delivery struct {
	SourceID    uint64
	MessageID   uint64
	MessageType string
	Encoding    string
	Timestamp   time.Time
	Payload     []byte
	Upstream    interface{}
	Receipt     interface{}
	Message     interface{}
}
