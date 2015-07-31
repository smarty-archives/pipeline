package rabbit

import (
	"errors"
	"net/url"
	"sync"
)

type Broker struct {
	mutex      *sync.Mutex
	target     url.URL
	connector  Connector
	connection Connection
	state      uint64
	readers    []Reader
	writers    []Writer
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

	if this.state == disconnecting {
		return ErrorShuttingDown
	} else if this.state == disconnected {
		this.state = connecting
	}

	return nil
}

func (this *Broker) Disconnect() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.state == disconnecting || this.state == disconnected {
		return
	}

	this.state = disconnecting

	if this.initiateReaderShutdown() {
		return // readers are busy shutting down
	}

	this.initiateWriterShutdown()
}
func (this *Broker) initiateReaderShutdown() bool {
	for _, reader := range this.readers {
		reader.Close()
	}

	return len(this.readers) > 0
}

func (this *Broker) initiateWriterShutdown() {
	for _, writer := range this.writers {
		writer.Close()
	}

	this.writers = this.writers[0:0]
	this.state = disconnected
	if this.connection != nil {
		this.connection.Close()
		this.connection = nil
	}
}

func (this *Broker) removeReader(reader interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for i, item := range this.readers {
		if reader != item {
			continue
		}

		this.readers = append(this.readers[:i], this.readers[i+1:]...)
		break
	}

	if this.state != disconnecting {
		return
	}

	if len(this.readers) == 0 {
		this.initiateWriterShutdown() // when all readers shutdown processes have been completed
	}
}

const (
	disconnected = iota
	connecting
	connected
	disconnecting
)

var ErrorShuttingDown = errors.New("Broker is still shutting down.")
