package rabbit

import (
	"errors"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/streadway/amqp"
)

type ChannelWriterFixture struct {
	*gunit.Fixture

	writer     *ChannelWriter
	controller *FakeWriterController
}

func (this *ChannelWriterFixture) Setup() {
	this.controller = newFakeWriterController()
	this.writer = newWriter(this.controller)
}

///////////////////////////////////////////////////////////////

func (this *ChannelWriterFixture) TestDispatchIsWrittenToChannel() {
	dispatch := Dispatch{
		Destination: "destination",
		Payload:     []byte{1, 2, 3, 4, 5},
	}

	err := this.writer.Write(dispatch)

	this.So(err, should.BeNil)
	this.So(this.controller.channel.exchange, should.Equal, dispatch.Destination)
	this.So(this.controller.channel.dispatch.Body, should.Resemble, dispatch.Payload)
}

///////////////////////////////////////////////////////////////

func (this *ChannelWriterFixture) TestChannelCannotBeObtained() {
	this.controller.channel = nil

	err := this.writer.Write(Dispatch{})

	this.So(err, should.NotBeNil)
}

///////////////////////////////////////////////////////////////

func (this *ChannelWriterFixture) TestFailedChannelClosed() {
	this.controller.channel.err = errors.New("channel failed")

	err := this.writer.Write(Dispatch{})

	this.So(err, should.Equal, this.controller.channel.err)
	this.So(this.controller.channel.closed, should.Equal, 1)
	this.So(this.writer.channel, should.BeNil)
}

///////////////////////////////////////////////////////////////

type FakeWriterController struct {
	channel        *FakeWriterChannel
	removedWriters []Writer
}

func newFakeWriterController() *FakeWriterController {
	return &FakeWriterController{channel: newFakeWriterChannel()}
}

func (this *FakeWriterController) openChannel() Channel {
	if this.channel == nil {
		return nil // interface quirks require this hack
	}

	return this.channel
}
func (this *FakeWriterController) removeReader(reader Reader) {}
func (this *FakeWriterController) removeWriter(writer Writer) {
	this.removedWriters = append(this.removedWriters, writer)
}

func (this *FakeWriterController) Dispose() {
	this.channel = nil
}

///////////////////////////////////////////////////////////////

type FakeWriterChannel struct {
	closed   int
	writes   int
	exchange string
	dispatch amqp.Publishing
	err      error
}

func newFakeWriterChannel() *FakeWriterChannel {
	return &FakeWriterChannel{}
}

func (this *FakeWriterChannel) ConfigureChannelBuffer(int) error                     { return nil }
func (this *FakeWriterChannel) DeclareTransientQueue() (string, error)               { return "", nil }
func (this *FakeWriterChannel) BindExchangeToQueue(string, string) error             { return nil }
func (this *FakeWriterChannel) Consume(string, string) (<-chan amqp.Delivery, error) { return nil, nil }
func (this *FakeWriterChannel) ExclusiveConsume(string, string) (<-chan amqp.Delivery, error) {
	return nil, nil
}
func (this *FakeWriterChannel) CancelConsumer(string) error { return nil }
func (this *FakeWriterChannel) Close() error {
	this.closed++
	return nil
}
func (this *FakeWriterChannel) AcknowledgeSingleMessage(uint64) error          { return nil }
func (this *FakeWriterChannel) AcknowledgeMultipleMessages(value uint64) error { return nil }
func (this *FakeWriterChannel) ConfigureChannelAsTransactional() error         { return nil }
func (this *FakeWriterChannel) PublishMessage(destination string, dispatch amqp.Publishing) error {
	this.exchange = destination
	this.dispatch = dispatch
	this.writes++
	return this.err
}
func (this *FakeWriterChannel) CommitTransaction() error {
	return this.err
}
func (this *FakeWriterChannel) RollbackTransaction() error { return nil }
