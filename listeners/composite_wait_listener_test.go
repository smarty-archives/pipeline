package listeners

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
	"github.com/smartystreets/gunit"
)

type CompositeWaitListenerFixture struct {
	*gunit.Fixture

	completed time.Time
	listener  *CompositeWaitListener
	items     []Listener
}

func (this *CompositeWaitListenerFixture) Setup() {
	this.items = []Listener{&FakeListener{}, &FakeListener{}}
	this.listener = NewCompositeWaitListener(this.items...)
}

//////////////////////////////////////////

func (this *CompositeWaitListenerFixture) TestAllListenersAreCalledAndWaitedFor() {
	this.listener.Listen()

	this.completed = clock.UTCNow()

	for _, item := range this.items {
		if item == nil {
			continue
		}

		this.So(item.(*FakeListener).calls, should.Equal, 1)
		this.So(this.completed.After(item.(*FakeListener).instant), should.BeTrue)
	}
}

//////////////////////////////////////////

func (this *CompositeWaitListenerFixture) TestNilListenersDontCausePanic() {
	this.listener = NewCompositeWaitListener(nil, nil, nil)
	this.So(this.listener.Listen, should.NotPanic)
	this.So(this.listener.Close, should.NotPanic)
}

//////////////////////////////////////////

func (this *CompositeWaitListenerFixture) TestCloseCallsInnerListeners() {
	this.listener.Close()

	for _, item := range this.items {
		this.So(item.(*FakeListener).closeCalls, should.Equal, 1)
	}
}

func (this *CompositeWaitListenerFixture) TestMultipleCloseCallInnerListenersExactlyOnce() {
	this.listener.Close()
	this.listener.Close()

	for _, item := range this.items {
		this.So(item.(*FakeListener).closeCalls, should.Equal, 1)
	}
}

func (this *CompositeWaitListenerFixture) TestCloseDoesntInvokeInfiniteLoop() {
	this.listener = NewCompositeWaitShutdownListener(this.items...)

	go this.listener.Close()
	this.listener.Listen()

	for _, item := range this.items {
		this.So(item.(*FakeListener).closeCalls, should.Equal, 1)
	}
}

//////////////////////////////////////////

type FakeListener struct {
	calls      int
	closeCalls int
	instant    time.Time
}

func (this *FakeListener) Listen() {
	this.instant = clock.UTCNow()
	time.Sleep(time.Millisecond)
	this.calls++
}

func (this *FakeListener) Close() {
	this.closeCalls++
}

//////////////////////////////////////////
