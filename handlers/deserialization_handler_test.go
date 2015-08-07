package pipeline

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/pipeline/messaging"
)

type DeserializationHandlerFixture struct {
	*gunit.Fixture

	input        chan messaging.Delivery
	output       chan messaging.Delivery
	deserializer *FakeDeserializer
	handler      *DeserializationHandler
}

func (this *DeserializationHandlerFixture) Setup() {
	this.input = make(chan messaging.Delivery, 2)
	this.output = make(chan messaging.Delivery, 2)
	this.deserializer = &FakeDeserializer{}
	this.handler = NewDeserializationHandler(this.input, this.output, this.deserializer)
}

func (this *DeserializationHandlerFixture) TestHandler() {
	in1 := messaging.Delivery{MessageID: 42}
	in2 := messaging.Delivery{MessageID: 43}
	this.input <- in1
	this.input <- in2
	close(this.input)

	this.handler.Handle()

	out1 := <-this.output
	out2 := <-this.output

	this.So(out1.MessageID, should.Equal, in1.MessageID)
	this.So(out1.Message, should.Equal, "Deserialized!")

	this.So(out2.MessageID, should.Equal, in2.MessageID)
	this.So(out2.Message, should.Equal, "Deserialized!")

	// The output channel should have been closed.
	counter := 0
	for range this.output {
		counter++
	}
	this.So(counter, should.Equal, 0)
}

////////////////////////////////////////////////////////////////////////

type FakeDeserializer struct{}

func (this *FakeDeserializer) Deserialize(delivery *messaging.Delivery) {
	delivery.Message = "Deserialized!"
}

////////////////////////////////////////////////////////////////////////
