package handlers

import "github.com/smartystreets/messaging"

type RequestContext interface {
	Write(interface{})
	Close()
}

type WaitGroup interface {
	Add(delta int)
	Done()
}

type MessageHandler interface {
	Handle(interface{}) interface{}
}

type ApplicationHandler interface {
	Handle(interface{})
}

type Transformer interface {
	Transform(*messaging.Delivery)
}

type Sender interface {
	Send(interface{}) interface{}
}
