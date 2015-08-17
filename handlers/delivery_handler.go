package handlers

import "github.com/smartystreets/pipeline/messaging"

type DeliveryHandler struct {
	input       <-chan messaging.Delivery
	output      chan<- interface{}
	writer      messaging.CommitWriter
	application MessageHandler
}

func NewDeliveryHandler(input <-chan messaging.Delivery,
	output chan<- interface{},
	writer messaging.CommitWriter,
	application MessageHandler) *DeliveryHandler {

	return &DeliveryHandler{
		input:       input,
		output:      output,
		writer:      writer,
		application: application,
	}
}

func (this *DeliveryHandler) Listen() {
	for delivery := range this.input {
		delivery = delivery
		converted := this.application.Handle(delivery.Message)
		this.write(converted)
		this.tryCommit(delivery.Receipt)
	}

	close(this.output)
}

func (this *DeliveryHandler) write(message interface{}) {
	if message == nil {
		return
	} else if multiple, ok := message.([]interface{}); ok {
		for _, item := range multiple {
			this.dispatch(item)
		}
	} else {
		this.dispatch(message)
	}

}
func (this *DeliveryHandler) dispatch(message interface{}) {
	dispatch := messaging.Dispatch{Message: message}
	this.writer.Write(dispatch)
}

func (this *DeliveryHandler) tryCommit(receipt interface{}) {
	if len(this.input) > 0 {
		return
	}

	this.writer.Commit()
	this.output <- receipt
}
