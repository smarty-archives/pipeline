package handlers

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/pipeline/messaging"
)

type SimpleDeliveryHandlerFixture struct {
	*gunit.Fixture

	input       chan messaging.Delivery
	output      chan interface{}
	handler     *SimpleDeliveryHandler
	application *FakeSimpleApplication
}

func (this *SimpleDeliveryHandlerFixture) Setup() {
	this.input = make(chan messaging.Delivery, 4)
	this.output = make(chan interface{}, 4)
	this.application = &FakeSimpleApplication{}
	this.handler = NewSimpleDeliveryHandler(this.application, this.input, this.output)
}

func (this *SimpleDeliveryHandlerFixture) TestDeliveriesReceivedAndAcknowledged() {
	this.input <- messaging.Delivery{Message: "a", Receipt: 1}
	this.input <- messaging.Delivery{Message: "b", Receipt: 2}
	this.input <- messaging.Delivery{Message: "c", Receipt: 3}
	close(this.input)

	this.handler.Listen()

	this.So(this.application.handled, should.Resemble, []interface{}{"a", "b", "c"})
	this.So(<-this.output, should.Equal, 1)
	this.So(<-this.output, should.Equal, 2)
	this.So(<-this.output, should.Equal, 3)
	this.So(<-this.output, should.BeNil)
}

type FakeSimpleApplication struct{ handled []interface{} }

func (this *FakeSimpleApplication) Handle(message interface{}) {
	this.handled = append(this.handled, message)
}
