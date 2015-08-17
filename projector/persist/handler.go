package persist

import (
	"sync"
	"time"

	"github.com/smartystreets/metrics"
	"github.com/smartystreets/pipeline/projector"
)

type Handler struct {
	input  chan projector.DocumentMessage
	output chan<- interface{}
	writer Writer

	pending map[string]projector.Document
	waiter  *sync.WaitGroup
}

func NewHandler(input chan projector.DocumentMessage, output chan<- interface{}, writer Writer) *Handler {
	return &Handler{
		input:   input,
		output:  output,
		writer:  writer,
		pending: make(map[string]projector.Document),
		waiter:  new(sync.WaitGroup),
	}
}

func (this *Handler) Listen() {
	for message := range this.input {
		metrics.Measure(depthPersistQueue, int64(len(this.input)))

		this.addToBatch(message)

		if len(this.input) == 0 {
			this.handleCurrentBatch(message.Receipt)
		}
	}

	close(this.output) // TODO: add test
}

func (this *Handler) addToBatch(message projector.DocumentMessage) {
	for _, document := range message.Documents {
		this.pending[document.Path()] = document
	}
}

func (this *Handler) handleCurrentBatch(receipt interface{}) {
	this.persistPendingDocuments()
	this.sendLatestAcknowledgement(receipt)
	this.prepareForNextBatch()
}

func (this *Handler) persistPendingDocuments() {
	this.waiter.Add(len(this.pending))
	metrics.Measure(documentsToSave, int64(len(this.pending)))

	for _, document := range this.pending {
		go this.persist(document)
	}

	this.waiter.Wait()
}
func (this *Handler) persist(document projector.Document) {
	started := time.Now()
	this.writer.Write(document)
	metrics.Measure(documentWriteLatency, milliseconds(time.Since(started)))
	this.waiter.Done()
}

func milliseconds(duration time.Duration) int64 { return microseconds(duration) / 1000 }
func microseconds(duration time.Duration) int64 { return int64(duration.Nanoseconds() / 1000) }

func (this *Handler) sendLatestAcknowledgement(receipt interface{}) {
	this.output <- receipt
}

func (this *Handler) prepareForNextBatch() {
	this.pending = make(map[string]projector.Document)
}

var (
	depthPersistQueue    = metrics.AddGauge("???:persist-phase-backlog-depth", time.Second)         // TODO: application-specific
	documentsToSave      = metrics.AddGauge("???:documents-to-save", time.Second)                   // TODO: application-specific
	documentWriteLatency = metrics.AddGauge("???:document-write-latency-milliseconds", time.Second) // TODO: application-specific
)
