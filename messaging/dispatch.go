package messaging

import "time"

type Dispatch struct {
	SourceID        uint64
	MessageID       uint64
	Destination     string
	MessageType     string
	ContentType     string
	ContentEncoding string
	Durable         bool
	Timestamp       time.Time
	Expiration      time.Duration
	Payload         []byte
	Message         interface{}
}
