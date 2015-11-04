package httpx

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/pipeline/handlers"
)

type MessageSenderFixture struct {
	*gunit.Fixture

	waiter  *FakeWaiter
	channel chan handlers.EventMessage
	sender  *MessageSender
}

func (this *MessageSenderFixture) Setup() {
	this.waiter = &FakeWaiter{}
	this.channel = make(chan handlers.EventMessage, 4)
	this.sender = NewMessageSender(this.channel, func() WaitGroup { return this.waiter })
}

func (this *MessageSenderFixture) TestMessageSent() {
	this.sender.Send(42)

	this.So((<-this.channel).Message, should.Equal, 42)
	this.So([]time.Time{this.waiter.addCalled, this.waiter.waitCalled}, should.BeChronological)
}
