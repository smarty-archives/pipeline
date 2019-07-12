package handlers

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestHTTPRequestContextFixture(t *testing.T) {
	gunit.Run(new(HTTPRequestContextFixture), t)
}

type HTTPRequestContextFixture struct {
	*gunit.Fixture
	waiter  *FakeWaiter
	context *HTTPRequestContext
}

func (this *HTTPRequestContextFixture) Setup() {
	this.waiter = &FakeWaiter{}
	this.context = NewRequestContext(this.waiter)
}

// ///////////////////////////////////////////////

func (this *HTTPRequestContextFixture) TestNewContext() {
	this.So(this.waiter.counter, should.Equal, 1)
	this.So(this.waiter.addCalls, should.Equal, 1)
	this.So(this.waiter.doneCalls, should.Equal, 0)
}

// ///////////////////////////////////////////////

func (this *HTTPRequestContextFixture) TestNoWritesToContext() {
	this.So(this.context.Written(), should.BeNil)
}

func (this *HTTPRequestContextFixture) TestWriteToContext() {
	this.context.Write(42)

	this.So(this.context.Written(), should.Equal, 42)

	// waiter has not changed since constructor
	this.So(this.waiter.counter, should.Equal, 1)
	this.So(this.waiter.addCalls, should.Equal, 1)
	this.So(this.waiter.doneCalls, should.Equal, 0)
}

func (this *HTTPRequestContextFixture) TestMultipleWritesToContext() {
	this.context.Write(42)
	this.context.Write(43)
	this.context.Write(44)

	this.So(this.context.Written(), should.Resemble, []interface{}{42, 43, 44})
}

// ///////////////////////////////////////////////

func (this *HTTPRequestContextFixture) TestCloseInvokesWaiter() {
	this.context.Close()

	this.So(this.waiter.counter, should.Equal, 0)
	this.So(this.waiter.addCalls, should.Equal, 1)
	this.So(this.waiter.doneCalls, should.Equal, 1)
}

// ///////////////////////////////////////////////

type FakeWaiter struct {
	addCalls, doneCalls, waitCalls, counter int
	addCalled, waitCalled                   time.Time
}

func (this *FakeWaiter) Add(delta int) {
	this.addCalled = time.Now()
	this.addCalls++
	this.counter += delta
}

func (this *FakeWaiter) Done() {
	this.doneCalls++
	this.counter--
}

func (this *FakeWaiter) Wait() {
	this.waitCalls++
	this.waitCalled = time.Now()
}
