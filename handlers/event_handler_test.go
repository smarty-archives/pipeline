package handlers

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/messaging/v2"
)

func TestEventHandlerFixture(t *testing.T) {
	gunit.Run(new(EventHandlerFixture), t)
}

type EventHandlerFixture struct {
	*gunit.Fixture

	input   chan messaging.Delivery
	output  chan interface{}
	router  *FakeDomain
	locker  *FakeLocker
	handler *EventHandler
}

func (this *EventHandlerFixture) Setup() {
	this.input = make(chan messaging.Delivery, 16)
	this.output = make(chan interface{}, 16)
	this.router = &FakeDomain{}
	this.locker = &FakeLocker{}
	this.handler = NewEventHandler(this.input, this.output, this.router, this.locker)
}

////////////////////////////////////////////////////////////

func (this *EventHandlerFixture) TestMessagePassedToDomain() {
	this.input <- messaging.Delivery{Message: 1, Receipt: 11}
	this.input <- messaging.Delivery{Message: 2, Receipt: 12}
	this.input <- messaging.Delivery{Message: 3, Receipt: 13}
	close(this.input)

	this.handler.Listen()

	this.So(this.locker.locks, should.Equal, 1)
	this.So(this.locker.unlocks, should.Equal, 1)
	this.So(this.router.received, should.Resemble, []interface{}{1, 2, 3})
	this.So(<-this.output, should.Equal, 13)
	this.So(<-this.output, should.BeNil) // outbound channel is closed when Listen exits
}

////////////////////////////////////////////////////////////

func (this *EventHandlerFixture) TestNilLockerDoesntPanic() {
	this.handler = NewEventHandler(this.input, this.output, this.router, nil)
	this.input <- messaging.Delivery{Message: 3, Receipt: 13}
	close(this.input)

	this.So(this.handler.Listen, should.NotPanic)
}

////////////////////////////////////////////////////////////

type FakeDomain struct {
	received []interface{}
}

func (this *FakeDomain) Apply(event interface{}) bool {
	this.received = append(this.received, event)
	return true
}
