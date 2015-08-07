package listeners

import (
	"fmt"
	"strings"
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
	this.inner = &FakeForWaitGroupListener{waiter: this.waiter}
	this.listener = NewWaitGroupListener(this.inner, this.waiter)
}

func (this *WaitGroupListenerFixture) TestWaitGroupListenerCallsDone() {
	this.listener.Listen()
	this.waiter.Wait() // ensures Done() is called
	this.So(this.inner.working, should.BeTrue)
}

////////////////////////////////////////////////////////////////////////////////

type FakeForWaitGroupListener struct {
	working bool
	waiter  *sync.WaitGroup
}

func (this *FakeForWaitGroupListener) Listen() {
	value := fmt.Sprintf("%+v", this.waiter)
	this.working = strings.Contains(value, "counter:1") // ensures .Add(1) is called
}

////////////////////////////////////////////////////////////////////////////////
