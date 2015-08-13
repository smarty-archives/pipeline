package rabbit

import (
	"strconv"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/pipeline/messaging"
	"github.com/streadway/amqp"
)

type SubscriptionFixture struct {
	*gunit.Fixture

	queue        string
	bindings     []string
	channel      *FakeSubscriptionChannel
	subscription *Subscription

	control chan interface{}
	output  chan messaging.Delivery
}

func (this *SubscriptionFixture) Setup() {
	this.channel = newFakeSubscriptionChannel()
	this.control = make(chan interface{}, 4)
	this.output = make(chan messaging.Delivery, 8)
}
func (this *SubscriptionFixture) createSubscription() {
	this.subscription = newSubscription(
		this.channel, this.queue, this.bindings,
		this.control, this.output)
}

//////////////////////////////////////////////////////////////////

func (this *SubscriptionFixture) TestQueuedBasedSubscription() {
	this.queue = "test-queue"

	this.assertListen()

	this.So(this.channel.queue, should.Equal, this.queue)
	this.So(this.channel.exclusive, should.BeFalse)
}
func (this *SubscriptionFixture) TestExclusiveSubscription() {
	this.bindings = []string{"exchange1", "exchange2"}

	this.assertListen()

	this.So(this.channel.exclusive, should.BeTrue)
	this.So(this.channel.queue, should.NotBeEmpty)
	this.So(this.channel.boundQueue[0], should.Equal, this.channel.queue)
}
func (this *SubscriptionFixture) TestFailingAMQPChannel() {
	this.queue = "test-queue"
	this.channel.incoming = nil

	this.assertListen()

	this.So(this.channel.queue, should.NotBeEmpty)
}
func (this *SubscriptionFixture) assertListen() {
	this.createSubscription()

	go this.subscription.Listen()
	this.channel.close()

	this.So(this.channel.bufferSize, should.Equal, cap(this.output))
	this.So(this.channel.bindings, should.Resemble, this.bindings)
	this.So(this.channel.consumer, should.NotBeEmpty)
	this.So((<-this.control).(subscriptionClosed).DeliveryCount, should.Equal, 0)
}

//////////////////////////////////////////////////////////////////

func (this *SubscriptionFixture) TestDeliveriesArePushedToTheApplication() {
	this.queue = "test-queue"
	delivery1 := amqp.Delivery{Type: "test-message", Body: []byte{1, 2, 3, 4, 5}}
	delivery2 := amqp.Delivery{Type: "test-message2", Body: []byte{6, 7, 8, 9, 10}}

	this.channel.incoming <- delivery1
	this.channel.incoming <- delivery2
	close(this.channel.incoming)

	this.createSubscription()
	go this.subscription.Listen()

	this.So((<-this.output), should.Resemble, messaging.Delivery{
		MessageType: "test-message",
		Payload:     []byte{1, 2, 3, 4, 5},
		Receipt:     newReceipt(this.channel, 0),
		Upstream:    delivery1,
	})
	this.So((<-this.output), should.Resemble, messaging.Delivery{
		MessageType: "test-message2",
		Payload:     []byte{6, 7, 8, 9, 10},
		Receipt:     newReceipt(this.channel, 0),
		Upstream:    delivery2,
	})
	this.So((<-this.control).(subscriptionClosed).DeliveryCount, should.Equal, 2)
}

//////////////////////////////////////////////////////////////////

func (this *SubscriptionFixture) TestConsumerCancellation() {
	this.createSubscription()
	this.subscription.Close()
	this.So(this.channel.cancelled, should.BeTrue)
	this.So(this.channel.consumer, should.NotBeEmpty)
}

//////////////////////////////////////////////////////////////////

type FakeSubscriptionChannel struct {
	bufferSize int
	queue      string
	consumer   string
	boundQueue []string
	bindings   []string
	exclusive  bool
	cancelled  bool
	incoming   chan amqp.Delivery
}

func newFakeSubscriptionChannel() *FakeSubscriptionChannel {
	return &FakeSubscriptionChannel{
		incoming: make(chan amqp.Delivery, 16),
	}
}

func (this *FakeSubscriptionChannel) ConfigureChannelBuffer(value int) error {
	this.bufferSize = value
	return nil
}
func (this *FakeSubscriptionChannel) DeclareTransientQueue() (string, error) {
	return strconv.FormatInt(time.Now().UnixNano(), 10), nil
}
func (this *FakeSubscriptionChannel) BindExchangeToQueue(queue string, exchange string) error {
	this.boundQueue = append(this.boundQueue, queue)
	this.bindings = append(this.bindings, exchange)
	return nil
}

func (this *FakeSubscriptionChannel) Consume(queue, consumer string) (<-chan amqp.Delivery, error) {
	this.queue = queue
	this.consumer = consumer
	return this.incoming, nil
}
func (this *FakeSubscriptionChannel) ExclusiveConsume(queue string, consumer string) (<-chan amqp.Delivery, error) {
	this.queue = queue
	this.consumer = consumer
	this.exclusive = true
	return this.incoming, nil
}

func (this *FakeSubscriptionChannel) CancelConsumer(consumer string) error {
	this.cancelled = true
	this.consumer = consumer
	return nil
}

func (this *FakeSubscriptionChannel) AcknowledgeSingleMessage(uint64) error    { return nil }
func (this *FakeSubscriptionChannel) AcknowledgeMultipleMessages(uint64) error { return nil }
func (this *FakeSubscriptionChannel) Close() error                             { return nil }

func (this *FakeSubscriptionChannel) close() {
	time.Sleep(time.Millisecond)
	if this.incoming != nil {
		close(this.incoming)
	}
}
