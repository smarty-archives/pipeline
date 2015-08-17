package projector

import "time"

type (
	RawMessage struct { // TODO: remove?
		Payload         []byte
		MessageType     string
		Acknowledgement DeliveryReceipt
	}

	TransformationMessage struct {
		Message         interface{}
		Acknowledgement DeliveryReceipt
		Now             time.Time
	}

	DocumentMessage struct {
		Documents []Document
		Receipt   interface{}
	}

	DeliveryReceipt interface {
		Acknowledge()
	}

	Document interface {
		Lapse(now time.Time) (next Document)
		Handle(message interface{}) bool
		Path() string
	}

	Listener interface {
		Listen()
	}
)
