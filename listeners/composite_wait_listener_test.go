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
	fakes     []*FakeListener
}

func (this *CompositeWaitListenerFixture) Setup() {
	this.fakes = []*FakeListener{&FakeListener{}, &FakeListener{}, nil}

	var fakes []Listener
	for _, fake := range this.fakes {
		fakes = append(fakes, fake)
	}

	this.listener = NewCompositeWaitListener(fakes...)
}

//////////////////////////////////////////

func (this *CompositeWaitListenerFixture) TestAllListenersAreCalledAndWaitedFor() {
	this.listener.Listen()

	this.completed = clock.UTCNow()

	for _, fake := range this.fakes {
		if fake == nil {
			continue
		}

		this.So(fake.calls, should.Equal, 1)
		this.So(this.completed.After(fake.instant), should.BeTrue)
	}
}

//////////////////////////////////////////

func (this *CompositeWaitListenerFixture) Test() {
}

//////////////////////////////////////////

type FakeListener struct {
	calls   int
	instant time.Time
}

func (this *FakeListener) Listen() {
	if this == nil {
		return
	}

	this.instant = clock.UTCNow()
	time.Sleep(time.Millisecond)
	this.calls++
}

//////////////////////////////////////////
