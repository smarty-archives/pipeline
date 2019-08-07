package handlers

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/messaging/v2"
)

func TestDeliveryHandlerFixture(t *testing.T) {
	gunit.Run(new(DeliveryHandlerFixture), t)
}

type DeliveryHandlerFixture struct {
	*gunit.Fixture

	now         time.Time
	input       chan messaging.Delivery
	output      chan interface{}
	writer      *FakeCommitWriter
	application *FakeApplication
	handler     *DeliveryHandler
	locker      *FakeLocker
}

func (this *DeliveryHandlerFixture) Setup() {
	this.input = make(chan messaging.Delivery, 8)
	this.output = make(chan interface{}, 8)
	this.writer = &FakeCommitWriter{}
	this.application = &FakeApplication{}
	this.locker = &FakeLocker{}
	this.handler = NewDeliveryHandler(this.input, this.output, this.writer, this.application, this.locker)
}

///////////////////////////////////////////////////////////////

func (this *DeliveryHandlerFixture) TestCommittCalledAtEndOfBatch() {
	this.input <- messaging.Delivery{Message: 1, Receipt: "Delivery Receipt 1"}
	this.input <- messaging.Delivery{Message: 2, Receipt: "Delivery Receipt 2"}
	this.input <- messaging.Delivery{Message: 3, Receipt: "Delivery Receipt 3"}

	close(this.input)
	this.handler.Listen()

	this.So(this.writer.commits, should.Equal, 1)
	this.So(len(this.output), should.Equal, 1)
	this.So(<-this.output, should.Equal, "Delivery Receipt 3")
	this.So(this.locker.locks, should.Equal, 1)
	this.So(this.locker.unlocks, should.Equal, 1)
}

///////////////////////////////////////////////////////////////

func (this *DeliveryHandlerFixture) TestNilMessagesAreNotDelivered() {
	this.input <- messaging.Delivery{Message: nil}
	close(this.input)

	this.handler.Listen()

	this.So(this.application.counter, should.BeZeroValue)
}

///////////////////////////////////////////////////////////////

func (this *DeliveryHandlerFixture) TestOutputChannelClosed() {
	close(this.input)
	this.handler.Listen()

	this.So(<-this.output, should.Equal, nil)
}

///////////////////////////////////////////////////////////////

func (this *DeliveryHandlerFixture) TestApplicationGeneratedMessagesAreWritten() {
	this.input <- messaging.Delivery{Message: 10, Receipt: "Delivery Receipt 1"}
	this.input <- messaging.Delivery{Message: 11, Receipt: "Delivery Receipt 2"}
	this.input <- messaging.Delivery{Message: 12, Receipt: "Delivery Receipt 3"}

	close(this.input)
	this.handler.Listen()

	this.So(this.writer.written, should.Resemble, []messaging.Dispatch{
		{Message: 1, Durable: true},
		{Message: 2, Durable: true},
		{Message: 3, Durable: true},
	})

	this.So(this.writer.commits, should.Equal, 1)
	this.So(len(this.output), should.Equal, 1)
	this.So(<-this.output, should.Equal, "Delivery Receipt 3")
}

///////////////////////////////////////////////////////////////

func (this *DeliveryHandlerFixture) TestNilMessagesAreNotWritten() {
	this.input <- messaging.Delivery{Message: "nil", Receipt: "Delivery Receipt"}

	close(this.input)
	this.handler.Listen()

	this.So(this.writer.written, should.BeEmpty)
	this.So(this.writer.commits, should.Equal, 1)
	this.So(len(this.output), should.Equal, 1)
	this.So(<-this.output, should.Equal, "Delivery Receipt")
	this.So(this.locker.locks, should.Equal, 1)
	this.So(this.locker.unlocks, should.Equal, 1)
}

///////////////////////////////////////////////////////////////

func (this *DeliveryHandlerFixture) TestMessageSlicesAreWritten() {
	this.input <- messaging.Delivery{Message: "multiple", Receipt: "Delivery Receipt"}

	close(this.input)
	this.handler.Listen()

	this.So(this.writer.written, should.Resemble, []messaging.Dispatch{
		{Message: 1, Durable: true},
		{Message: 2, Durable: true},
		{Message: 3, Durable: true},
	})
	this.So(this.writer.commits, should.Equal, 1)
	this.So(len(this.output), should.Equal, 1)
	this.So(<-this.output, should.Equal, "Delivery Receipt")
	this.So(this.locker.locks, should.Equal, 1)
	this.So(this.locker.unlocks, should.Equal, 1)
}

///////////////////////////////////////////////////////////////

type FakeCommitWriter struct {
	written []messaging.Dispatch
	commits int
}

func (this *FakeCommitWriter) Write(dispatch messaging.Dispatch) error {
	this.written = append(this.written, dispatch)
	return nil
}

func (this *FakeCommitWriter) Commit() error {
	this.commits++
	return nil
}

func (this *FakeCommitWriter) Close() {
	panic("should never be called")
}

///////////////////////////////////////////////////////////////

type FakeApplication struct {
	counter int
}

func (this *FakeApplication) Handle(message interface{}) interface{} {
	if message == "nil" {
		return nil
	} else if message == "multiple" {
		return []interface{}{1, 2, 3}
	}

	this.counter++
	return this.counter
}

///////////////////////////////////////////////////////////////
