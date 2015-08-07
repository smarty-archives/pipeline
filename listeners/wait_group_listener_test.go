package listeners

import (
	"sync"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type WaitGroupListenerFixture struct {
	*gunit.Fixture

	inner    *FakeForWaitGroupListener
	waiter   *sync.WaitGroup
	listener Listener
}

func (this *WaitGroupListenerFixture) Setup() {
	this.waiter = &sync.WaitGroup{}
	this.inner = &FakeForWaitGroupListener{}
	this.listener = NewWaitGroupListener(this.inner, this.waiter)
}

func (this *WaitGroupListenerFixture) TestWaitGroupListenerCallsDone() {
	go this.listener.Listen()

	this.waiter.Wait()
	this.So(this.inner.called, should.Equal, 1)
}

////////////////////////////////////////////////////////////////////////////////

type FakeForWaitGroupListener struct {
	called int
}

func (this *FakeForWaitGroupListener) Listen() {
	this.called++
}

////////////////////////////////////////////////////////////////////////////////
