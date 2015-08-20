package rabbit

import (
	"sync"

	"github.com/smartystreets/pipeline/messaging"
)

type ChannelReader struct {
	mutex             *sync.Mutex
	controller        Controller
	queue             string
	bindings          []string
	control           chan interface{}
	acknowledgements  chan interface{}
	deliveries        chan messaging.Delivery
	shutdown          bool
	shutdownRequested bool
}

func newReader(controller Controller, queue string, bindings []string) *ChannelReader {
	return &ChannelReader{
		mutex:            &sync.Mutex{},
		controller:       controller,
		queue:            queue,
		bindings:         bindings,
		control:          make(chan interface{}, 32),
		acknowledgements: make(chan interface{}, 1024),
		deliveries:       make(chan messaging.Delivery, 1024),
	}
}

func (this *ChannelReader) Listen() {
	acknowledger := newAcknowledger(this.control, this.acknowledgements)
	go acknowledger.Listen()

	for this.listen() {
	}

	close(this.deliveries)
	this.controller.removeReader(this)
}
func (this *ChannelReader) listen() bool {
	channel := this.controller.openChannel(this.isActive)
	if channel == nil {
		return false // broker no longer allowed to give me a channel, it has been manually closed
	}

	subscription := this.subscribe(channel)

	for element := range this.control {
		switch item := element.(type) {
		case shutdownRequested:
			this.shutdown = true
			subscription.Close()
		case subscriptionClosed:
			if this.shutdown {
				// keep channel alive and gracefully stop acknowledgement listener
				this.acknowledgements <- item
			} else {
				// channel failure; reconnect
				channel.Close()
				return true
			}
		case acknowledgementCompleted:
			channel.Close() // we don't need the channel anymore
			return false    // the shutdown process for this reader is complete
		}
	}

	return true
}
func (this *ChannelReader) subscribe(channel Channel) *Subscription {
	subscription := newSubscription(channel, this.queue, this.bindings, this.control, this.deliveries)
	go subscription.Listen()
	return subscription
}

func (this *ChannelReader) Close() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if !this.shutdownRequested {
		this.control <- shutdownRequested{}
		this.shutdownRequested = true
	}
}
func (this *ChannelReader) isActive() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return !this.shutdownRequested
}

func (this *ChannelReader) Deliveries() <-chan messaging.Delivery {
	return this.deliveries
}
func (this *ChannelReader) Acknowledgements() chan<- interface{} {
	return this.acknowledgements
}
