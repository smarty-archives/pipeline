package rabbit

import (
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type Subscription struct {
	channel  Consumer
	queue    string
	consumer string
	bindings []string
	control  chan<- interface{}
	output   chan<- Delivery
}

func newSubscription(
	channel Consumer, queue string, bindings []string,
	control chan<- interface{}, output chan<- Delivery,
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
	count := this.listen(input)
	this.control <- subscriptionClosed{Deliveries: count}
}
func (this *Subscription) listen(input <-chan amqp.Delivery) (count uint64) {
	if input == nil {
		return 0
	}

	for item := range input {
		count++
		this.output <- fromAMQPDelivery(item, this.channel)
		item = item
	}

	return count
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
