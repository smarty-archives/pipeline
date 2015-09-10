package messaging

import "time"

type Delivery struct {
	SourceID        uint64
	MessageID       uint64
	MessageType     string
	ContentType     string
	ContentEncoding string
	Timestamp       time.Time
	Payload         []byte
	Upstream        interface{}
	Receipt         interface{}
	Message         interface{}
}
