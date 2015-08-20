package rabbit

import (
	"strconv"
	"time"

	"github.com/smartystreets/pipeline/messaging"
	"github.com/streadway/amqp"
)

type Subscription struct {
	channel      Consumer
	queue        string
	consumer     string
	bindings     []string
	finalReceipt interface{}
	control      chan<- interface{}
	output       chan<- messaging.Delivery
}

func newSubscription(
	channel Consumer, queue string, bindings []string,
	control chan<- interface{}, output chan<- messaging.Delivery,
) *Subscription {
	return &Subscription{
		channel:  channel,
		queue:    queue,
		consumer: strconv.FormatInt(time.Now().UnixNano(), 10),
		bindings: bindings,
		control:  control,
		output:   output,
	}
}

func (this *Subscription) Listen() {
	input := this.open()
	this.listen(input)
	this.control <- subscriptionClosed{FinalReceipt: this.finalReceipt}
}
func (this *Subscription) listen(input <-chan amqp.Delivery) {
	if input == nil {
		return
	}

	for item := range input {
		delivery := fromAMQPDelivery(item, this.channel)
		this.finalReceipt = delivery.Receipt
		this.output <- delivery
	}
}
func (this *Subscription) open() <-chan amqp.Delivery {
	this.channel.ConfigureChannelBuffer(cap(this.output))

	if len(this.queue) > 0 {
		return this.consume()
	}

	this.queue, _ = this.channel.DeclareTransientQueue()
	for _, exchange := range this.bindings {
		this.channel.BindExchangeToQueue(this.queue, exchange)
	}

	return this.exclusiveConsume()
}

func (this *Subscription) consume() <-chan amqp.Delivery {
	queue, _ := this.channel.Consume(this.queue, this.consumer)
	return queue
}
func (this *Subscription) exclusiveConsume() <-chan amqp.Delivery {
	queue, _ := this.channel.ExclusiveConsume(this.queue, this.consumer)
	return queue
}

func (this *Subscription) Close() {
	this.channel.CancelConsumer(this.consumer)
}
