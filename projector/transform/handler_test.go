package transform

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/pipeline/projector"
)

type HandlerFixture struct {
	*gunit.Fixture

	input       chan projector.TransformationMessage
	output      chan projector.DocumentMessage
	transformer *FakeTransformer
	handler     *Handler
	firstInput  projector.TransformationMessage
	secondInput projector.TransformationMessage
}

func (this *HandlerFixture) Setup() {
	this.input = make(chan projector.TransformationMessage, 2)
	this.output = make(chan projector.DocumentMessage, 2)
	this.transformer = NewFakeTransformer()
	this.handler = NewHandler(this.input, this.output, this.transformer)

	this.firstInput = projector.TransformationMessage{
		Message:         1,
		Now:             time.Now(),
		Acknowledgement: &FakeAcknowledgement{},
	}
	this.secondInput = projector.TransformationMessage{
		Message:         2,
		Now:             time.Now(),
		Acknowledgement: &FakeAcknowledgement{},
	}
}

/////////////////////////////////////////////////////////////////

func (this *HandlerFixture) TestTransformerInvokedForEveryInputMessage() {
	this.input <- this.firstInput
	this.input <- this.secondInput
	close(this.input)

	this.handler.Listen()

	this.So(this.transformer.received, should.Resemble, map[interface{}]time.Time{
		this.firstInput.Message:  this.firstInput.Now,
		this.secondInput.Message: this.secondInput.Now,
	})
	this.So(<-this.output, should.Resemble, projector.DocumentMessage{
		Acknowledgement: this.secondInput.Acknowledgement,
		Documents:       collectedDocuments,
	})
}

/////////////////////////////////////////////////////////////////

type FakeTransformer struct {
	received map[interface{}]time.Time
}

func NewFakeTransformer() *FakeTransformer {
	return &FakeTransformer{
		received: make(map[interface{}]time.Time),
	}
}

func (this *FakeTransformer) TransformAllDocuments(message interface{}, now time.Time) {
	this.received[message] = now
}

var collectedDocuments = []projector.Document{
	&fakeDocument{path: "a"},
	&fakeDocument{path: "b"},
	&fakeDocument{path: "c"},
}

func (this *FakeTransformer) Collect() []projector.Document {
	return collectedDocuments
}

/////////////////////////////////////////////////////////////////

type FakeAcknowledgement struct{}

func (this *FakeAcknowledgement) Acknowledge() {}

/////////////////////////////////////////////////////////////////

type fakeDocument struct{ path string }

func (this *fakeDocument) Path() string                                  { return this.path }
func (this *fakeDocument) Lapse(now time.Time) (next projector.Document) { return this }
func (this *fakeDocument) Handle(message interface{}) bool               { return false }

/////////////////////////////////////////////////////////////////
