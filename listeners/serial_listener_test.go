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
		&FakeForSerialListener{listened: time.Now().Add(time.Second)},
		&FakeForSerialListener{listened: time.Now()},
		&FakeForSerialListener{listened: time.Now().Add(-time.Second)},
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
}

type FakeForSerialListener struct {
	calls    int
	listened time.Time
}

func (this *FakeForSerialListener) Listen() {
	this.calls++
	this.listened = clock.UTCNow()
	time.Sleep(time.Microsecond)
}
