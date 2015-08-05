package rabbit

import (
	"net/url"
	"sync"
	"time"

	"github.com/smartystreets/clock"
	"github.com/smartystreets/go-messenger"
)

type Broker struct {
	mutex      *sync.Mutex
	target     url.URL
	connector  Connector
	connection Connection
	readers    []messenger.Reader
	writers    []messenger.Writer
	state      uint64
}

func NewBroker(target url.URL, connector Connector) *Broker {
	return &Broker{
		mutex:     &sync.Mutex{},
		target:    target,
		connector: connector,
	}
}

func (this *Broker) State() uint64 {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.state
}

func (this *Broker) Connect() error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.state == messenger.Disconnecting {
		return messenger.BrokerShuttingDownError
	} else if this.state == messenger.Disconnected {
		this.state = messenger.Connecting
	}

	return nil
}

func (this *Broker) Disconnect() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.state == messenger.Disconnecting || this.state == messenger.Disconnected {
		return
	}

	this.state = messenger.Disconnecting

	this.initiateReaderShutdown()
	this.initiateWriterShutdown()
	this.completeShutdown()
}
func (this *Broker) initiateReaderShutdown() {
	for _, reader := range this.readers {
		reader.Close()
	}
}

func (this *Broker) initiateWriterShutdown() {
	if len(this.readers) > 0 {
		return
	}

	for _, writer := range this.writers {
		writer.Close()
	}

	this.writers = this.writers[0:0]
}
func (this *Broker) completeShutdown() {
	if this.state != messenger.Disconnecting {
		return
	}

	if len(this.readers) > 0 || len(this.writers) > 0 {
		return
	}

	if this.connection != nil {
		this.connection.Close()
		this.connection = nil
	}

	this.state = messenger.Disconnected
}

func (this *Broker) removeReader(reader messenger.Reader) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for i, item := range this.readers {
		if reader != item {
			continue
		}

		this.readers = append(this.readers[:i], this.readers[i+1:]...)
		break
	}

	if this.state != messenger.Disconnecting {
		return
	}

	this.initiateWriterShutdown() // when all readers shutdown processes have been completed
	this.completeShutdown()
}
func (this *Broker) removeWriter(writer messenger.Writer) {
	for i, item := range this.writers {
		if writer != item {
			continue
		}

		this.writers = append(this.writers[:i], this.writers[i+1:]...)
		break
	}
}

func (this *Broker) OpenReader(queue string) messenger.Reader {
	return this.openReader(queue, nil)
}
func (this *Broker) OpenTransientReader(bindings []string) messenger.Reader {
	return this.openReader("", bindings)
}
func (this *Broker) openReader(queue string, bindings []string) messenger.Reader {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.state == messenger.Disconnecting || this.state == messenger.Disconnected {
		return nil
	}

	reader := newReader(this, queue, bindings)
	this.readers = append(this.readers, reader)
	return reader
}

func (this *Broker) OpenWriter() messenger.Writer {
	writer := this.openWriter(false)
	this.writers = append(this.writers, writer)
	return writer
}
func (this *Broker) OpenTransactionalWriter() messenger.CommitWriter {
	writer := this.openWriter(true).(messenger.CommitWriter)
	this.writers = append(this.writers, writer)
	return writer
}
func (this *Broker) openWriter(transactional bool) messenger.Writer {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.state == messenger.Disconnecting || this.state == messenger.Disconnected {
		return nil
	}

	if transactional {
		return transactionWriter(this)
	} else {
		return newWriter(this)
	}
}

func (this *Broker) openChannel() Channel {
	// don't lock for the duration of the loop
	// otherwise this can deadlock because we're dependent
	// upon the broker to be online/active. By avoiding
	// a lock here, we can try to connect and if that fails
	// we can still shutdown properly

	for this.isActive() {
		if channel := this.tryOpenChannel(); channel != nil {
			return channel
		}

		clock.Sleep(time.Second * 4)
	}

	return nil
}
func (this *Broker) isActive() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.state == messenger.Connecting || this.state == messenger.Connected
}
func (this *Broker) tryOpenChannel() Channel {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if !this.ensureConnection() {
		return nil
	}

	return this.openChannelFromExistingConnection()
}
func (this *Broker) ensureConnection() bool {
	if this.connection != nil {
		return true
	}

	var err error
	this.connection, err = this.connector.Connect(this.target)
	return err == nil
}
func (this *Broker) openChannelFromExistingConnection() Channel {
	// remember to only change the state (this.connection, this.state)
	// within the protection of this.mutex
	if channel, err := this.connection.Channel(); err != nil {
		this.connection.Close()
		this.connection = nil
		this.state = messenger.Connecting
		return nil
	} else {
		this.state = messenger.Connected
		return channel
	}
}
