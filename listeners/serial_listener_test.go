package listeners

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
	"github.com/smartystreets/gunit"
)

type SerialListenerFixture struct {
	*gunit.Fixture
}

func (this *SerialListenerFixture) TestListenCallInOrder() {
	items := []Listener{
		&FakeForSerialListener{listened: clock.UTCNow().Add(time.Second)},
		&FakeForSerialListener{listened: clock.UTCNow()},
		&FakeForSerialListener{listened: clock.UTCNow().Add(-time.Second)},
	}

	NewSerialListener(items...).Listen()

	times := []time.Time{}
	for _, item := range items {
		fake := item.(*FakeForSerialListener)
		times = append(times, fake.listened)

		this.So(fake.calls, should.Equal, 1)
	}

	this.So(times, should.BeChronological)

}

func (this *SerialListenerFixture) TestNilListenersAreIgnored() {
	this.So(NewSerialListener(nil).Listen, should.NotPanic)
	this.So(NewSerialListener(nil).Close, should.NotPanic)
}

func (this *SerialListenerFixture) TestCloseCalledOnInnerListeners() {
	items := []Listener{&FakeForSerialListener{}, &FakeForSerialListener{}}

	NewSerialListener(items...).Close()

	for _, item := range items {
		this.So(item.(*FakeForSerialListener).closeCalls, should.Equal, 1)
	}
}

type FakeForSerialListener struct {
	calls      int
	closeCalls int
	listened   time.Time
}

func (this *FakeForSerialListener) Listen() {
	this.calls++
	this.listened = clock.UTCNow()
	time.Sleep(time.Microsecond)
}

func (this *FakeForSerialListener) Close() {
	this.closeCalls++
}
