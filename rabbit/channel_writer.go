package rabbit

import (
	"errors"
	"sync"

	"github.com/smartystreets/clock"
	"github.com/smartystreets/go-messenger"
)

type ChannelWriter struct {
	mutex           *sync.Mutex
	controller      Controller
	channel         Channel
	closed          bool
	skipUntilCommit bool
}

func newWriter(controller Controller) *ChannelWriter {
	return &ChannelWriter{mutex: &sync.Mutex{}, controller: controller}
}

func (this *ChannelWriter) Write(message messenger.Dispatch) error {
	if !this.ensureChannel() {
		return channelFailure
	}

	dispatch := toAMQPDispatch(message, clock.Now())
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

	if this.closed {
		return false
	}

	this.channel = this.controller.openChannel()
	return this.channel != nil
}

var channelFailure = errors.New("Unable to obtain a connection and channel to the broker.")
