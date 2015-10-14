package rabbit

import (
	"log"
	"sync"

	"github.com/smartystreets/clock"
	"github.com/smartystreets/pipeline/messaging"
)

type TransactionWriter struct {
	mutex      *sync.Mutex
	controller Controller
	channel    Channel
	closed     bool
}

func transactionWriter(controller Controller) *TransactionWriter {
	return &TransactionWriter{
		mutex:      &sync.Mutex{},
		controller: controller,
	}
}

func (this *TransactionWriter) Write(message messaging.Dispatch) error {
	if !this.ensureChannel() {
		return messaging.WriterClosedError
	}

	// FUTURE: if error on publish, don't publish anything else
	// until we reset the channel during commit
	// opening a new channel is what marks it as able to continue
	dispatch := toAMQPDispatch(message, clock.UTCNow())
	return this.channel.PublishMessage(message.Destination, dispatch)
}

func (this *TransactionWriter) Commit() error {
	if this.channel == nil {
		return nil
	}

	err := this.channel.CommitTransaction()
	if err == nil {
		return nil
	}

	log.Println("[WARN] Transaction failed, closing channel: [", err, "]")
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

	this.channel = this.controller.openChannel(this.isActive)
	if this.channel == nil {
		return false
	}

	this.channel.ConfigureChannelAsTransactional()
	return true
}

func (this *TransactionWriter) isActive() bool {
	return !this.closed // must be called from within the safety of a mutex
}
