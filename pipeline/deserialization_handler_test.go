package pipeline

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/go-messenger"
	"github.com/smartystreets/gunit"
)

type DeserializationHandlerFixture struct {
	*gunit.Fixture

	input        chan messenger.Delivery
	output       chan messenger.Delivery
	deserializer *FakeDeserializer
	handler      *DeserializationHandler
}

func (this *DeserializationHandlerFixture) Setup() {
	this.input = make(chan messenger.Delivery, 2)
	this.output = make(chan messenger.Delivery, 2)
	this.deserializer = &FakeDeserializer{}
	this.handler = NewDeserializationHandler(this.input, this.output, this.deserializer)
}

func (this *DeserializationHandlerFixture) TestHandler() {
	in1 := messenger.Delivery{MessageID: 42}
	in2 := messenger.Delivery{MessageID: 43}
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

func (this *FakeDeserializer) Deserialize(delivery *messenger.Delivery) {
	delivery.Message = "Deserialized!"
}

////////////////////////////////////////////////////////////////////////
