package listeners

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type AckListenerFixture struct {
	*gunit.Fixture

	listener *AckListener
	input    chan interface{}
}

func (this *AckListenerFixture) Setup() {
	this.input = make(chan interface{}, 3)
	this.listener = NewAckListener(this.input)
}

func (this *AckListenerFixture) TestWaitGroupsThatArriveAreMarkedAsDone__UnknownItemsAreSilentlyIgnored() {
	waiter1 := &EmptyWaitGroup{}
	waiter2 := &EmptyWaitGroup{}

	this.input <- time.Now() // Definitely not a waitgroup, will be ignored.
	this.input <- waiter1
	this.input <- waiter2
	close(this.input)

	this.listener.Listen()

	this.So(waiter1.done, should.BeTrue)
	this.So(waiter2.done, should.BeTrue)
}

////////////////////////////////////////////////////////////////////////////////

type EmptyWaitGroup struct{ done bool }

func (this *EmptyWaitGroup) Add(delta int) {}
func (this *EmptyWaitGroup) Done()         { this.done = true }
