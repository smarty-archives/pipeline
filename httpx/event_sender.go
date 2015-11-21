package httpx

import "github.com/smartystreets/pipeline/handlers"

type EventSender struct {
	output  chan<- handlers.EventMessage
	factory func() WaitGroup
}

func NewEventSender(output chan<- handlers.EventMessage, factory func() WaitGroup) *EventSender {
	return &EventSender{
		output:  output,
		factory: factory,
	}
}

func (this *EventSender) Send(message interface{}) {
	waiter := this.factory()
	this.output <- handlers.EventMessage{
		Message:    message,
		Context:    NewRequestContext(waiter),
		EndOfBatch: true,
	}
	waiter.Wait()
}
