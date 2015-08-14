package handlers

import (
	"sync"

	"github.com/smartystreets/pipeline/domain"
)

type CommandHandler struct {
	input  <-chan RequestMessage
	output chan<- EventMessage
	router domain.Handler
	locker sync.Locker
}

func NewCommandHandler(
	input <-chan RequestMessage,
	output chan<- EventMessage,
	router domain.Handler,
	locker sync.Locker) *CommandHandler {

	return &CommandHandler{
		input:  input,
		output: output,
		router: router,
		locker: NewIdempotentLocker(locker),
	}
}

func (this *CommandHandler) Listen() {
	for item := range this.input {
		this.processCommand(item)
	}

	close(this.output)
}

func (this *CommandHandler) processCommand(item RequestMessage) {
	this.locker.Lock()

	events := this.router.Handle(item.Message)
	this.sendResultingEvents(item, events)

	if len(this.input) == 0 {
		this.locker.Unlock()
	}
}

func (this *CommandHandler) sendResultingEvents(item RequestMessage, events []interface{}) {
	if len(events) == 0 {
		this.sendEmpty(item.Context)
	} else {
		this.sendResults(item.Context, events)
	}
}

func (this *CommandHandler) sendEmpty(context RequestContext) {
	this.output <- EventMessage{Context: context, EndOfBatch: true}
}

func (this *CommandHandler) sendResults(context RequestContext, events []interface{}) {
	for i, event := range events {
		this.output <- EventMessage{
			Message:    event,
			Context:    context,
			EndOfBatch: len(events)-1 == i,
		}
	}
}
