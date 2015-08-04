package rabbit

import (
	"sync"

	"github.com/smartystreets/clock"
)

type TransactionWriter struct {
	mutex           *sync.Mutex
	controller      Controller
	channel         Channel
	closed          bool
	skipUntilCommit bool
}

func transactionWriter(controller Controller) *TransactionWriter {
	return &TransactionWriter{
		mutex:      &sync.Mutex{},
		controller: controller,
	}
}

func (this *TransactionWriter) Write(message Dispatch) error {
	if !this.ensureChannel() {
		return channelFailure
	}

	dispatch := toAMQPDispatch(message, clock.Now())
	return this.channel.PublishMessage(message.Destination, dispatch)
}

func (this *TransactionWriter) Commit() error {
	if this.channel == nil {
		return channelFailure // this never creates a channel
	}

	err := this.channel.CommitTransaction()
	if err == nil {
		return nil
	}

	this.channel.Close()
	this.channel = nil
	return err
}

func (this *TransactionWriter) Close() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.closed = true
}

func (this *TransactionWriter) ensureChannel() bool {
	if this.channel != nil {
		return true
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return false
	}

	this.channel = this.controller.openChannel()
	if this.channel == nil {
		return false
	}

	this.channel.ConfigureChannelAsTransactional()
	return true
}
