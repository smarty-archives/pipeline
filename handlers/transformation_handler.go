package handlers

import "github.com/smartystreets/messaging/v2"

type TransformationHandler struct {
	input        <-chan messaging.Delivery
	output       chan<- messaging.Delivery
	transformers []Transformer
}

func NewTransformationHandler(
	input <-chan messaging.Delivery,
	output chan<- messaging.Delivery,
	transformers ...Transformer) *TransformationHandler {

	return &TransformationHandler{
		input:        input,
		output:       output,
		transformers: transformers,
	}
}

func (this *TransformationHandler) Listen() {
	for delivery := range this.input {
		this.transform(&delivery)
		this.output <- delivery
	}

	close(this.output)
}

func (this *TransformationHandler) transform(delivery *messaging.Delivery) {
	for _, transformer := range this.transformers {
		if transformer == nil {
			continue
		}

		transformer.Transform(delivery)
	}
}
