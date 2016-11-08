package handlers

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestContextAcknowledgementHandlerFixture(t *testing.T) {
	gunit.Run(new(ContextAcknowledgementHandlerFixture), t)
}

type ContextAcknowledgementHandlerFixture struct {
	*gunit.Fixture

	listener *ContextAcknowledgementHandler
	input    chan RequestContext
}

func (this *ContextAcknowledgementHandlerFixture) Setup() {
	this.input = make(chan RequestContext, 4)
	this.listener = NewContextAcknowledgementHandler(this.input)
}

func (this *ContextAcknowledgementHandlerFixture) TestContextsThatArriveClosed() {
	waiter1 := &EmptyWaitGroup{}
	waiter2 := &EmptyWaitGroup{}
	this.input <- &FakeRequestContext{waiter: waiter1}
	this.input <- &FakeRequestContext{waiter: waiter2}
	close(this.input)

	this.listener.Listen()

	this.So(waiter1.done, should.BeTrue)
	this.So(waiter2.done, should.BeTrue)
}

////////////////////////////////////////////////////////////////////////////////

type FakeRequestContext struct {
	waiter  WaitGroup
	written []interface{}
	id      int
}

func (this *FakeRequestContext) Write(item interface{}) {
	this.written = append(this.written, item)
}
func (this *FakeRequestContext) Close() { this.waiter.Done() }

type EmptyWaitGroup struct{ done bool }

func (this *EmptyWaitGroup) Add(delta int) {}
func (this *EmptyWaitGroup) Done()         { this.done = true }
