package listeners

import (
	"sync"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type WaitGroupListenerFixture struct {
	*gunit.Fixture

	inner    *FakeForWaitGroupListener
	waiter   *WrappedWaitGroup
	listener *WaitGroupListener
}

func (this *WaitGroupListenerFixture) Setup() {
	this.waiter = NewWrappedWaitGroup()
	this.inner = &FakeForWaitGroupListener{}
	this.listener = NewWaitGroupListener(this.inner, this.waiter)
}

func (this *WaitGroupListenerFixture) TestWaitGroupListenerCallsDone() {
	this.listener.Listen()
	this.waiter.Wait() // This ensures that .Add(1) and .Done() were called.
	this.So(this.waiter.added, should.BeTrue)
	this.So(this.waiter.done, should.BeTrue)
	this.So(this.inner.listenCalls, should.Equal, 1)
}

////////////////////////////////////////////////////////////////////////////////

func (this *WaitGroupListenerFixture) TestNilInnerListener() {
	this.listener = NewWaitGroupListener(nil, this.waiter)

	this.So(this.listener.Listen, should.NotPanic)
	this.So(this.listener.Close, should.NotPanic)
}

////////////////////////////////////////////////////////////////////////////////

func (this *WaitGroupListenerFixture) TestCloseCallsInner() {
	this.listener.Close()

	this.So(this.inner.closeCalls, should.Equal, 1)
}

////////////////////////////////////////////////////////////////////////////////

type FakeForWaitGroupListener struct {
	listenCalls int
	closeCalls  int
}

func (this *FakeForWaitGroupListener) Listen() {
	this.listenCalls++
}

func (this *FakeForWaitGroupListener) Close() {
	this.closeCalls++
}

////////////////////////////////////////////////////////////////////////////////

type WrappedWaitGroup struct {
	inner *sync.WaitGroup
	added bool
	done  bool
}

func NewWrappedWaitGroup() *WrappedWaitGroup {
	return &WrappedWaitGroup{inner: new(sync.WaitGroup)}
}

func (this *WrappedWaitGroup) Add(delta int) {
	this.added = true
	this.inner.Add(delta)
}

func (this *WrappedWaitGroup) Done() {
	this.done = true
	this.inner.Done()
}

func (this *WrappedWaitGroup) Wait() {
	this.inner.Wait()
}

////////////////////////////////////////////////////////////////////////////////
