package handlers

import "sync"

type RequestSender struct {
	output chan<- RequestMessage
}

func NewRequestSender(output chan<- RequestMessage) *RequestSender {
	return &RequestSender{output: output}
}

func (this *RequestSender) Send(message interface{}) interface{} {
	waiter := new(sync.WaitGroup)
	context := NewRequestContext(waiter)
	this.output <- RequestMessage{
		Message: message,
		Context: context,
	}
	waiter.Wait()
	return context.Written()
}
