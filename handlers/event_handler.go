package handlers

import (
	"sync"

	"github.com/smartystreets/pipeline/domain"
	"github.com/smartystreets/pipeline/messaging"
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

		if len(this.input) == 0 {
			this.locker.Unlock()
			this.output <- delivery.Receipt
		}
	}

	close(this.output)
}
