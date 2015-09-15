package transform

import (
	"time"

	"github.com/smartystreets/clock"
	"github.com/smartystreets/metrics"
	"github.com/smartystreets/pipeline/messaging"
	"github.com/smartystreets/pipeline/projector"
)

type Handler struct {
	input       <-chan messaging.Delivery
	output      chan<- projector.DocumentMessage
	transformer Transformer
	clock       *clock.Clock
}

func NewHandler(input <-chan messaging.Delivery, output chan<- projector.DocumentMessage, transformer Transformer) *Handler {
	return &Handler{input: input, output: output, transformer: transformer}
}

func (this *Handler) Listen() {
	for message := range this.input {
		now := this.clock.UTCNow()

		metrics.Measure(transformQueueDepth, int64(len(this.input)))

		this.transformer.TransformAllDocuments(message.Message, now)

		if len(this.input) == 0 {
			this.output <- projector.DocumentMessage{
				Receipt:   message.Receipt,
				Documents: this.transformer.Collect(),
			}
		}
	}

	close(this.output) // TODO: add test
}

var transformQueueDepth = metrics.AddGauge("pipeline:transform-phase-backlog-depth", time.Second)
