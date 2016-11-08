package rabbit

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/pipeline/messaging"
	"github.com/streadway/amqp"
)

func TestChannelReaderFixture(t *testing.T) {
	gunit.Run(new(ChannelReaderFixture), t)
}

type ChannelReaderFixture struct {
	*gunit.Fixture

	reader     *ChannelReader
	controller *FakeReaderController
	queue      string
	bindings   []string
}

func (this *ChannelReaderFixture) Setup() {
	this.queue = "queue"
	this.controller = newFakeReaderController()
	this.buildReader()
}
func (this *ChannelReaderFixture) buildReader() {
	this.reader = newReader(this.controller, this.queue, this.bindings)
}
func (this *ChannelReaderFixture) controlMessage() interface{} {
	return <-this.reader.control
}

///////////////////////////////////////////////////////////////

func (this *ChannelReaderFixture) TestCloseReader() {
	this.reader.Close()
	this.So(this.controlMessage(), should.Resemble, shutdownRequested{})
}

///////////////////////////////////////////////////////////////

func (this *ChannelReaderFixture) TestDisconnectedControllerExitsListenLoop() {
	this.controller.Dispose()

	this.reader.Listen()

	this.So(this.controller.removedReaders[0], should.Equal, this.reader)
}

///////////////////////////////////////////////////////////////

func (this *ChannelReaderFixture) TestListenStartsAcknowledger() {
	channel := newFakeReaderChannel()
	receipt := newReceipt(channel, 42)
	go this.reader.Listen()

	this.reader.Acknowledgements() <- receipt
	time.Sleep(time.Millisecond * 10)

	this.So(channel.latestTag, should.Equal, 42)
}

///////////////////////////////////////////////////////////////

func (this *ChannelReaderFixture) TestCloseShutsdownReaderAfterAllMessagesProcessed() {
	channel := this.controller.channel
	channel.deliveries <- amqp.Delivery{}

	go func() {
		message := <-this.reader.Deliveries()
		this.reader.Close()
		this.reader.Acknowledgements() <- message.Receipt
	}()

	this.reader.Listen()

	for range this.reader.Deliveries() {
	}

	this.So(true, should.BeTrue) // we only get here when the golang channel is closed
	this.So(channel.closed, should.Equal, 1)
}

///////////////////////////////////////////////////////////////

func (this *ChannelReaderFixture) TestCloseShutsdownReaderAfterLastMessageProcessed() {
	channel := this.controller.channel
	channel.deliveries <- amqp.Delivery{}
	channel.deliveries <- amqp.Delivery{DeliveryTag: 42}

	go func() {
		<-this.reader.Deliveries()
		last := <-this.reader.Deliveries()
		this.reader.Close()
		this.reader.Acknowledgements() <- last.Receipt
	}()

	this.reader.Listen()

	for range this.reader.Deliveries() {
	}

	this.So(true, should.BeTrue) // we only get here when the golang channel is closed
	this.So(channel.closed, should.Equal, 1)
}

///////////////////////////////////////////////////////////////

func (this *ChannelReaderFixture) TestDeliveriesChannelClosedWhenReaderCompleted() {
	this.controller.channel = nil
	this.reader.Listen()

	for range this.reader.Deliveries() {
	}

	this.So(true, should.BeTrue) // we only get here when the golang channel is closed
}

///////////////////////////////////////////////////////////////

func (this *ChannelReaderFixture) TestCloseOnlyClosesOnce() {
	this.reader.Close()
	this.reader.Close()

	this.So(len(this.reader.control), should.Equal, 1)
}

///////////////////////////////////////////////////////////////

type FakeReaderController struct {
	channel        *FakeReaderChannel
	removedReaders []messaging.Reader
}

func newFakeReaderController() *FakeReaderController {
	return &FakeReaderController{channel: newFakeReaderChannel()}
}

func (this *FakeReaderController) openChannel(callback func() bool) Channel {
	if !callback() {
		return nil
	}

	if this.channel == nil {
		return nil // interface quirks require this hack
	}

	return this.channel
}
func (this *FakeReaderController) removeReader(reader messaging.Reader) {
	this.removedReaders = append(this.removedReaders, reader)
}
func (this *FakeReaderController) removeWriter(writer messaging.Writer) {}

func (this *FakeReaderController) Dispose() {
	this.channel = nil
}

///////////////////////////////////////////////////////////////

type FakeReaderChannel struct {
	latestTag     uint64
	deliveries    chan amqp.Delivery
	closed        int
	cancellations int
}

func newFakeReaderChannel() *FakeReaderChannel {
	return &FakeReaderChannel{deliveries: make(chan amqp.Delivery, 16)}
}

func (this *FakeReaderChannel) ConfigureChannelBuffer(int) error         { return nil }
func (this *FakeReaderChannel) DeclareExchange(string, string) error     { return nil }
func (this *FakeReaderChannel) DeclareQueue(string) error                { return nil }
func (this *FakeReaderChannel) DeclareTransientQueue() (string, error)   { return "", nil }
func (this *FakeReaderChannel) BindExchangeToQueue(string, string) error { return nil }
func (this *FakeReaderChannel) Consume(string, string) (<-chan amqp.Delivery, error) {
	return this.deliveries, nil
}
func (this *FakeReaderChannel) ExclusiveConsume(string, string) (<-chan amqp.Delivery, error) {
	return this.deliveries, nil
}
func (this *FakeReaderChannel) ConsumeWithoutAcknowledgement(string, string) (<-chan amqp.Delivery, error) {
	return nil, nil
}
func (this *FakeReaderChannel) ExclusiveConsumeWithoutAcknowledgement(string, string) (<-chan amqp.Delivery, error) {
	return nil, nil
}
func (this *FakeReaderChannel) CancelConsumer(string) error {
	this.cancellations++
	close(this.deliveries)
	return nil
}
func (this *FakeReaderChannel) Close() error {
	this.closed++
	return nil
}
func (this *FakeReaderChannel) AcknowledgeSingleMessage(uint64) error { return nil }
func (this *FakeReaderChannel) AcknowledgeMultipleMessages(value uint64) error {
	this.latestTag = value
	return nil
}
func (this *FakeReaderChannel) ConfigureChannelAsTransactional() error       { return nil }
func (this *FakeReaderChannel) PublishMessage(string, amqp.Publishing) error { return nil }
func (this *FakeReaderChannel) CommitTransaction() error                     { return nil }
func (this *FakeReaderChannel) RollbackTransaction() error                   { return nil }
