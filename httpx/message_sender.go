package httpx

import "github.com/smartystreets/pipeline/handlers"

type MessageSender struct {
	output  chan<- handlers.EventMessage
	factory func() WaitGroup
}

func NewMessageSender(output chan<- handlers.EventMessage, factory func() WaitGroup) *MessageSender {
	return &MessageSender{
		output:  output,
		factory: factory,
	}
}

func (this *MessageSender) Send(message interface{}) {
	waiter := this.factory()
	this.output <- handlers.EventMessage{
		Message:    message,
		Context:    NewRequestContext(waiter),
		EndOfBatch: true,
	}
	waiter.Wait()
}
