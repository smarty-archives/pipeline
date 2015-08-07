package messaging

type Delivery struct {
	SourceID    uint64
	MessageID   uint64
	MessageType string
	Encoding    string
	Payload     []byte
	Upstream    interface{}
	Receipt     interface{}
	Message     interface{}
}
