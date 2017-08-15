package httpx

import "github.com/smartystreets/pipeline/handlers"

type RequestSender struct {
	output chan<- handlers.RequestMessage
}

func NewRequestSender(output chan<- handlers.RequestMessage) *RequestSender {
	return &RequestSender{output: output}
}

func (this *RequestSender) Send(message interface{}) interface{} {
	waiter := NewWaitGroup(0)
	context := NewRequestContext(waiter)
	this.output <- handlers.RequestMessage{
		Message: message,
		Context: context,
	}
	waiter.Wait()
	return context.Written()
}
