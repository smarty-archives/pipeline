package listeners

import (
	"os"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"
)

func TestShutdownListenerFixture(t *testing.T) {
	gunit.Run(new(ShutdownListenerFixture), t)
}

type ShutdownListenerFixture struct {
	*gunit.Fixture

	calls    int
	listener *ShutdownListener
}

func (this *ShutdownListenerFixture) Setup() {
	this.listener = NewShutdownListener(func() { this.calls++ })
	this.listener.logger = logging.Capture()
}

func (this *ShutdownListenerFixture) TestShutdownSignalInvokesShutdownCallback() {
	this.listener.channel <- os.Interrupt
	this.listener.Listen()
	this.So(this.calls, should.Equal, 1)
}

func (this *ShutdownListenerFixture) TestClosingBlockedListenerInvokesShutdownCallback() {
	go this.listener.Close()
	this.listener.Listen()
	this.So(this.calls, should.Equal, 1)
}

func (this *ShutdownListenerFixture) TestCloseBehaviorHappensOnce() {
	this.listener.Close()
	this.So(this.listener.Close, should.NotPanic)
}
