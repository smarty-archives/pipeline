package handlers

import "github.com/smartystreets/messaging/v2"

type SimpleDeliveryHandler struct {
	application ApplicationHandler
	input       <-chan messaging.Delivery
	output      chan<- interface{}
}

func NewSimpleDeliveryHandler(application ApplicationHandler,
	input <-chan messaging.Delivery,
	output chan<- interface{}) *SimpleDeliveryHandler {

	return &SimpleDeliveryHandler{
		application: application,
		input:       input,
		output:      output,
	}
}

func (this *SimpleDeliveryHandler) Listen() {
	for delivery := range this.input {
		this.application.Handle(delivery.Message)
		this.output <- delivery.Receipt
	}

	close(this.output)
}
