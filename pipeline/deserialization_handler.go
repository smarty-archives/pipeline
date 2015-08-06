package pipeline

import "github.com/smartystreets/go-messenger"

type DeserializationHandler struct {
	input        chan messenger.Delivery
	output       chan messenger.Delivery
	deserializer Deserializer
}

func NewDeserializationHandler(input, output chan messenger.Delivery, deserializer Deserializer) *DeserializationHandler {
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
