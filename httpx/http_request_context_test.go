package httpx

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type HTTPRequestContextFixture struct {
	*gunit.Fixture
	waiter  *FakeWaiter
	context *HTTPRequestContext
}

func (this *HTTPRequestContextFixture) Setup() {
	this.waiter = &FakeWaiter{}
	this.context = NewRequestContext(this.waiter)
}

/////////////////////////////////////////////////

func (this *HTTPRequestContextFixture) TestNewContext() {
	this.So(this.waiter.counter, should.Equal, 1)
	this.So(this.waiter.addCalls, should.Equal, 1)
	this.So(this.waiter.doneCalls, should.Equal, 0)
}

/////////////////////////////////////////////////

func (this *HTTPRequestContextFixture) TestWriteToContext() {
	this.context.Write(42)

	this.So(this.context.Result, should.Equal, 42)

	// waiter has not changed since constructor
	this.So(this.waiter.counter, should.Equal, 1)
	this.So(this.waiter.addCalls, should.Equal, 1)
	this.So(this.waiter.doneCalls, should.Equal, 0)
}

/////////////////////////////////////////////////

func (this *HTTPRequestContextFixture) TestCloseInvokesWaiter() {
	this.context.Close()

	this.So(this.waiter.counter, should.Equal, 0)
	this.So(this.waiter.addCalls, should.Equal, 1)
	this.So(this.waiter.doneCalls, should.Equal, 1)
}

/////////////////////////////////////////////////

type FakeWaiter struct{ addCalls, doneCalls, counter int }

func (this *FakeWaiter) Add(delta int) {
	this.addCalls++
	this.counter += delta
}

func (this *FakeWaiter) Done() {
	this.doneCalls++
	this.counter--
}
