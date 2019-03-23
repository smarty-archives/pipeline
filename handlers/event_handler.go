package handlers

import (
	"sync"

	"github.com/smartystreets/messaging"
	"github.com/smartystreets/pipeline/domain"
)

type EventHandler struct {
	input  <-chan messaging.Delivery
	output chan<- interface{}
	router domain.Applicator
	locker sync.Locker
}

func NewEventHandler(
	input <-chan messaging.Delivery,
	output chan<- interface{},
	router domain.Applicator,
	locker sync.Locker) *EventHandler {

	if locker == nil {
		locker = NoopLocker{}
	} else {
		locker = NewIdempotentLocker(locker)
	}

	return &EventHandler{
		input:  input,
		output: output,
		router: router,
		locker: locker,
	}
}

func (this *EventHandler) Listen() {
	for delivery := range this.input {
		this.locker.Lock()
		this.router.Apply(delivery.Message)

		// TODO: if the backlog is large, this could stop other pipelines from forward progress
		// we need to figure out how to do this gracefully
		if len(this.input) == 0 {
			this.locker.Unlock()
			this.output <- delivery.Receipt
		}
	}

	close(this.output)
}
