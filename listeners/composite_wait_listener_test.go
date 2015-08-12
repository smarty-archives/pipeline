package listeners

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type CompositeWaitListenerFixture struct {
	*gunit.Fixture

	listener *CompositeWaitListener
	fakes    []*FakeListener
}

func (this *CompositeWaitListenerFixture) Setup() {
	this.fakes = []*FakeListener{&FakeListener{}, &FakeListener{}, &FakeListener{}}

	var fakes []Listener
	for _, fake := range this.fakes {
		fakes = append(fakes, fake)
	}

	this.listener = NewCompositeWaitListener(fakes...)
}

//////////////////////////////////////////

func (this *CompositeWaitListenerFixture) TestAllListenersAreCalledAndWaitedFor() {
	this.listener.Listen()

	for _, fake := range this.fakes {
		this.So(fake.calls, should.Equal, 1)
	}
}

//////////////////////////////////////////

type FakeListener struct{ calls int }

func (this *FakeListener) Listen() {
	time.Sleep(time.Millisecond)
	this.calls++
}

//////////////////////////////////////////
