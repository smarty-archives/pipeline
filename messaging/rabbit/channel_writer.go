package rabbit

import (
	"sync"

	"github.com/smartystreets/clock"
	"github.com/smartystreets/pipeline/messaging"
)

type ChannelWriter struct {
	mutex      *sync.Mutex
	controller Controller
	channel    Channel
	closed     bool
}

func newWriter(controller Controller) *ChannelWriter {
	return &ChannelWriter{mutex: &sync.Mutex{}, controller: controller}
}

func (this *ChannelWriter) Write(message messaging.Dispatch) error {
	if !this.ensureChannel() {
		return messaging.WriterClosedError
	}

	dispatch := toAMQPDispatch(message, clock.UTCNow())
	err := this.channel.PublishMessage(message.Destination, dispatch)
	if err == nil {
		return nil
	}

	this.channel.Close()
	this.channel = nil
	return err
}

func (this *ChannelWriter) Commit() error {
	return nil
}

func (this *ChannelWriter) Close() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.closed = true
}

func (this *ChannelWriter) ensureChannel() bool {
	if this.channel != nil {
		return true
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.channel = this.controller.openChannel(this.isActive)
	return this.channel != nil
}
func (this *ChannelWriter) isActive() bool {
	return !this.closed // must be called within the safety of a mutex
}
