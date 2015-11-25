package persist

import (
	"strconv"
	"sync"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/pipeline/projector"
)

type HandlerFixture struct {
	*gunit.Fixture

	input   chan projector.DocumentMessage
	output  chan interface{}
	writer  *FakeWriter
	handler *Handler
}

func (this *HandlerFixture) Setup() {
	this.input = make(chan projector.DocumentMessage, 2)
	this.output = make(chan interface{}, 2)
	this.writer = NewFakeWriter()
	this.handler = NewHandler(this.input, this.output, this.writer)
}

func (this *HandlerFixture) send(messages ...projector.DocumentMessage) {
	for _, message := range messages {
		this.input <- message
	}
}
func (this *HandlerFixture) sendAndClose(messages ...projector.DocumentMessage) {
	this.send(messages...)
	close(this.input)
}
func (this *HandlerFixture) assertRecordsWritten(howMany int) bool {
	return this.So(len(this.writer.written), should.Equal, howMany)
}
func (this *HandlerFixture) assertDocumentsWritten(paths ...string) {
	for _, path := range paths {
		this.So(path, should.BeIn, this.writer.written)
	}
}
func (this *HandlerFixture) assertDocumentNotWritten(path string) {
	this.So(path, should.NotBeIn, this.writer.written)
}

//////////////////////////////////////////////////////////////

func (this *HandlerFixture) TestAllDocumentsWritten() {
	this.sendAndClose(newMessage(0, 0))
	this.handler.Listen()
	if this.assertRecordsWritten(2) {
		this.assertDocumentsWritten("0", "1")
	}
}

//////////////////////////////////////////////////////////////

func (this *HandlerFixture) TestAllDocumentsWrittenInAConsolidatedBatch() {
	finalAck := 100
	this.sendAndClose(newMessage(0, 0), newMessage(1, finalAck))
	this.handler.Listen()
	if this.assertRecordsWritten(3) {
		this.assertDocumentsWritten("0", "1", "2")
	}
	this.assertLatestDeliveryReceiptSentToNextPhase(finalAck)
	this.assertOutputChannelClosed()
}
func (this *HandlerFixture) assertLatestDeliveryReceiptSentToNextPhase(expectedAck int) {
	deliveryReceipt := (<-this.output).(*FakeReceipt)
	this.So(deliveryReceipt.id, should.Equal, expectedAck)
}
func (this *HandlerFixture) assertOutputChannelClosed() {
	for range this.output {
		// if the channel isn't closed, it will block here
	}
}

//////////////////////////////////////////////////////////////

func (this *HandlerFixture) TestEachSentBatchIsSeparateFromThePreviousOne() {
	go this.handler.Listen()

	this.Print("Send the first batch and clear out the writer")
	this.send(newMessage(0, 0))
	time.Sleep(time.Millisecond * 10) // allow the handler to consume and write the batch
	this.writer.written = []string{}

	this.Print("Send the second batch")
	this.send(newMessage(1, 1))
	time.Sleep(time.Millisecond * 10) // allow the handler to consume and write the batch
	close(this.input)

	this.Print("The first batch should not end up being written twice (with the second batch)")
	this.assertDocumentNotWritten("0") // it was written in a previous batch
	this.assertDocumentsWritten("1", "2")
}

//////////////////////////////////////////////////////////////

func newMessage(documentIndex, ack int) projector.DocumentMessage {
	return projector.DocumentMessage{
		Receipt: &FakeReceipt{id: ack},
		Documents: []projector.Document{
			NewFakeDocument(strconv.Itoa(documentIndex)),
			NewFakeDocument(strconv.Itoa(documentIndex + 1)),
		},
	}
}

//////////////////////////////////////////////////////////////

type FakeWriter struct {
	mutex   *sync.Mutex
	written []string
}

func NewFakeWriter() *FakeWriter {
	return &FakeWriter{mutex: &sync.Mutex{}}
}

func (this *FakeWriter) Write(document projector.Document) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.written = append(this.written, document.Path())
}

/////////////////////////////////////

type FakeDocument struct{ id string }

func NewFakeDocument(id string) *FakeDocument                 { return &FakeDocument{id: id} }
func (this *FakeDocument) Path() string                       { return this.id }
func (this *FakeDocument) Lapse(time.Time) projector.Document { panic("NOT IMPLEMENTED") }
func (this *FakeDocument) Apply(interface{}) bool             { panic("NOT IMPLEMENTED") }

/////////////////////////

type FakeReceipt struct{ id int }

func (this *FakeReceipt) Acknowledge() {}

/////////////////////////////////////////
