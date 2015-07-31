package rabbit

import (
	"github.com/smartystreets/gunit"
	"github.com/streadway/amqp"
)

type ChannelWriterFixture struct {
	*gunit.Fixture

	writer        *ChannelWriter
	controller    *FakeWriterController
	transactional bool
}

func (this *ChannelWriterFixture) Setup() {
	this.controller = newFakeWriterController()
	this.buildWriter()
}
func (this *ChannelWriterFixture) buildWriter() {
	this.writer = newWriter(this.controller, this.transactional)
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
	closed int
	writes int
	err    error
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
func (this *FakeWriterChannel) PublishMessage(string, amqp.Publishing) error {
	this.writes++
	return this.err
}
func (this *FakeWriterChannel) CommitTransaction() error {
	return this.err
}
func (this *FakeWriterChannel) RollbackTransaction() error { return nil }
