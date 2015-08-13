package rabbit

type BatchAcknowledger struct {
	control          chan<- interface{}
	input            <-chan interface{}
	count            uint64
	maximum          uint64
	finalTag         uint64
	finalConsumer    interface{}
	receivedTag      uint64
	receivedConsumer interface{}
	closing          bool
	pending          *DeliveryReceipt
}

func newAcknowledger(control chan<- interface{}, input <-chan interface{}) *BatchAcknowledger {
	return &BatchAcknowledger{control: control, input: input}
}

func (this *BatchAcknowledger) Listen() {
	this.listen()
	this.acknowledge()
	this.control <- acknowledgementCompleted{Acknowledgements: this.count}
}

func (this *BatchAcknowledger) listen() {
	for item := range this.input {
		this.processItem(item)
		if this.isComplete() {
			return
		}
	}
}

func (this *BatchAcknowledger) processItem(entity interface{}) {
	switch item := entity.(type) {
	case subscriptionClosed:
		this.processClosingEvent(item)
	case DeliveryReceipt:
		this.processAcknowledgment(item)
	}
}
func (this *BatchAcknowledger) processClosingEvent(item subscriptionClosed) {
	this.closing = true
	this.maximum += item.DeliveryCount
	this.finalTag = item.LatestDeliveryTag
	this.finalConsumer = item.LatestConsumer
}
func (this *BatchAcknowledger) processAcknowledgment(item DeliveryReceipt) {
	this.count++
	this.receivedTag = item.deliveryTag
	this.receivedConsumer = item.channel

	if len(this.input) > 0 {
		this.pending = &item
	} else {
		this.pending = nil
		acknowledge(item)
	}
}

func (this *BatchAcknowledger) acknowledge() {
	if this.pending != nil {
		acknowledge(*this.pending)
		this.pending = nil
	}
}
func acknowledge(receipt DeliveryReceipt) {
	receipt.channel.AcknowledgeMultipleMessages(receipt.deliveryTag)
}

func (this *BatchAcknowledger) isComplete() bool {
	if !this.closing {
		return false
	}

	if len(this.input) > 0 {
		return false
	}

	if this.maximum <= this.count {
		return true
	}

	if this.finalTag != this.receivedTag {
		return false
	}

	if this.finalConsumer != this.receivedConsumer {
		return false
	}

	return true
}
