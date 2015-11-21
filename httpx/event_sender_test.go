package httpx

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/pipeline/handlers"
)

type EventSenderFixture struct {
	*gunit.Fixture

	waiter  *FakeWaiter
	channel chan handlers.EventMessage
	sender  *EventSender
}

func (this *EventSenderFixture) Setup() {
	this.waiter = &FakeWaiter{}
	this.channel = make(chan handlers.EventMessage, 4)
	this.sender = NewEventSender(this.channel, func() WaitGroup { return this.waiter })
}

func (this *EventSenderFixture) TestMessageSent() {
	this.sender.Send(42)

	this.So((<-this.channel).Message, should.Equal, 42)
	this.So([]time.Time{this.waiter.addCalled, this.waiter.waitCalled}, should.BeChronological)
}
