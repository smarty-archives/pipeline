package messenger

type Delivery struct {
	SourceID    uint64
	MessageID   uint64
	MessageType string
	Encoding    string
	Payload     []byte
	Receipt     interface{}
	Message     interface{}
}
