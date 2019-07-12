package handlers

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestRequestSenderFixture(t *testing.T) {
	gunit.Run(new(RequestSenderFixture), t)
}

type RequestSenderFixture struct {
	*gunit.Fixture

	waiter  *FakeWaiter
	channel chan RequestMessage
	sender  *RequestSender
}

func (this *RequestSenderFixture) Setup() {
	this.waiter = &FakeWaiter{}
	this.channel = make(chan RequestMessage, 4)
	this.sender = NewRequestSender(this.channel)
}

func (this *RequestSenderFixture) TestMessageSent() {
	go func() {
		for item := range this.channel {
			if item.Message == "Hello," {
				item.Context.Write("World!")
				item.Context.Close()
			}
		}
	}()

	result := this.sender.Send("Hello,")

	this.So(result, should.Equal, "World!")
}
