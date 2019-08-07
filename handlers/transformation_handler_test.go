package handlers

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/messaging/v2"
)

func TestTransformationHandlerFixture(t *testing.T) {
	gunit.Run(new(TransformationHandlerFixture), t)
}

type TransformationHandlerFixture struct {
	*gunit.Fixture

	input   chan messaging.Delivery
	output  chan messaging.Delivery
	handler *TransformationHandler
}

func (this *TransformationHandlerFixture) Setup() {
	this.input = make(chan messaging.Delivery, 2)
	this.output = make(chan messaging.Delivery, 2)
	this.handler = NewTransformationHandler(this.input, this.output)
}

func (this *TransformationHandlerFixture) TestNoTransformers() {
	in1 := messaging.Delivery{MessageID: 42}
	in2 := messaging.Delivery{MessageID: 43}
	this.input <- in1
	this.input <- in2
	close(this.input)

	this.handler.Listen()

	this.So(<-this.output, should.Resemble, in1)
	this.So(<-this.output, should.Resemble, in2)
}

func (this *TransformationHandlerFixture) TestCloseOutputWhenInputClosed() {
	close(this.input)

	this.handler.Listen()

	counter := 0
	for i := range this.output {
		i = i
		counter++ // The output channel should have been closed (unreachable code)
	}
	this.So(counter, should.Equal, 0)
}

func (this *TransformationHandlerFixture) TestAllTransformersCalled() {
	this.handler = NewTransformationHandler(this.input, this.output,
		&FakeTransformer{}, &FakeTransformer{}, &FakeTransformer{})

	delivery := messaging.Delivery{MessageID: 42}
	this.input <- delivery
	close(this.input)

	this.handler.Listen()

	this.So((<-this.output).MessageID, should.Equal, 45)
}

type FakeTransformer struct {
	delivery *messaging.Delivery
}

func (this *FakeTransformer) Transform(delivery *messaging.Delivery) {
	delivery.MessageID++
	this.delivery = delivery
}
