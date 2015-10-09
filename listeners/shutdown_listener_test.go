package listeners

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type ShutdownListenerFixture struct {
	*gunit.Fixture
}

func (this *ShutdownListenerFixture) Setup() {
	log.SetOutput(ioutil.Discard)
}
func (this *ShutdownListenerFixture) Teardown() {
	log.SetOutput(os.Stdout)
}

func (this *ShutdownListenerFixture) TestShutdownSignalInvokesShutdownCallback() {
	var calls int
	listener := NewShutdownListener(func() { calls++ })
	listener.channel <- os.Interrupt

	listener.Listen()

	this.So(calls, should.Equal, 1)
}

func (this *ShutdownListenerFixture) TestClosingBlockedListenerInvokesShutdownCallback() {
	var calls int
	listener := NewShutdownListener(func() { calls++ })

	go listener.Close()
	listener.Listen()

	this.So(calls, should.Equal, 1)
}

func (this *ShutdownListenerFixture) TestCloseBehaviorHappensOnce() {
	listener := NewShutdownListener(func() {})

	listener.Close()

	this.So(listener.Close, should.NotPanic)
}
