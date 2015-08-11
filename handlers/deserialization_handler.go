package handlers

import "github.com/smartystreets/pipeline/messaging"

type DeserializationHandler struct {
	input        chan messaging.Delivery
	output       chan messaging.Delivery
	deserializer Deserializer
}

func NewDeserializationHandler(input, output chan messaging.Delivery, deserializer Deserializer) *DeserializationHandler {
	return &DeserializationHandler{
		input:        input,
		output:       output,
		deserializer: deserializer,
	}
}

func (this *DeserializationHandler) Handle() {
	for delivery := range this.input {
		this.deserializer.Deserialize(&delivery)
		this.output <- delivery
	}
	close(this.output)
}
