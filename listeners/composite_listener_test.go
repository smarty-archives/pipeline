package listeners

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type CompositeListenerFixture struct {
	*gunit.Fixture

	listeners []Listener
	composite Listener
}

func (this *CompositeListenerFixture) Setup() {
	for x := 0; x < 100; x++ {
		this.listeners = append(this.listeners, &FakeForCompositeListener{})
	}
	this.composite = NewCompositeListener(this.listeners...)
}

func (this *CompositeListenerFixture) TestCompositeListenerCallsInnerListenersConcurrently() {
	started := time.Now()
	this.composite.Listen()
	this.So(time.Since(started), should.BeLessThan, nap*2)
}

////////////////////////////////////////////////////////////////////////////////

func (this *CompositeListenerFixture) TestCompositeListenerDoesntFailWithNoListeners() {
	this.listeners = nil
	this.composite = NewCompositeListener(this.listeners...)
	this.So(this.composite.Listen, should.NotPanic)
}

////////////////////////////////////////////////////////////////////////////////

func (this *CompositeListenerFixture) TestCompositeListenerSkipNilListeners() {
	this.listeners = append(this.listeners, &FakeForCompositeListener{})
	this.listeners = append(this.listeners, nil)
	this.listeners = append(this.listeners, nil)
	this.composite = NewCompositeListener(this.listeners...)
	this.So(this.composite.Listen, should.NotPanic)
}

////////////////////////////////////////////////////////////////////////////////

var nap = time.Millisecond

type FakeForCompositeListener struct{}

func (this *FakeForCompositeListener) Listen() {
	time.Sleep(nap)
}

////////////////////////////////////////////////////////////////////////////////
