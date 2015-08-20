package rabbit

type BatchAcknowledger struct {
	control chan<- interface{}
	input   <-chan interface{}
	pending *DeliveryReceipt
	closing bool
	final   bool
	waiting interface{}
}

func newAcknowledger(control chan<- interface{}, input <-chan interface{}) *BatchAcknowledger {
	return &BatchAcknowledger{control: control, input: input}
}

func (this *BatchAcknowledger) Listen() {
	this.listen()
	this.acknowledge()
	this.control <- acknowledgementCompleted{}
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
	this.waiting = item.FinalReceipt
}
func (this *BatchAcknowledger) processAcknowledgment(item DeliveryReceipt) {
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

	if this.pending != nil {
		return false
	}

	if len(this.input) > 0 {
		return false
	}

	return this.final
}
