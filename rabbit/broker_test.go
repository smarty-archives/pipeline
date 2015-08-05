package rabbit

import (
	"errors"
	"net/url"
	"reflect"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
	"github.com/smartystreets/go-messenger"
	"github.com/smartystreets/gunit"
	"github.com/streadway/amqp"
)

type BrokerFixture struct {
	*gunit.Fixture

	target    url.URL
	connector *FakeConnector
	broker    *Broker

	sleeper *clock.Sleeper
}

func (this *BrokerFixture) Setup() {
	target, _ := url.Parse("amqps://username:password@localhost:5671/")
	this.target = *target
	this.connector = NewFakeConnector(0, 0)
	this.createBroker()

	this.sleeper = clock.FakeSleep()
}
func (this *BrokerFixture) Teardown() {
	this.sleeper.Restore()
}

func (this *BrokerFixture) createBroker() {
	this.broker = NewBroker(this.target, this.connector)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestConnect() {
	this.assertConnectResult(disconnected, connecting, false)
	this.assertConnectResult(connecting, connecting, false)
	this.assertConnectResult(connected, connected, false)
	this.assertConnectResult(disconnecting, disconnecting, true)
}
func (this *BrokerFixture) assertConnectResult(initial, updated uint64, hasError bool) {
	this.broker.state = initial

	err := this.broker.Connect()
	if hasError {
		this.So(err, should.NotBeNil)
	} else {
		this.So(err, should.BeNil)
	}

	this.So(this.broker.State(), should.Equal, updated)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestDisconnectWithoutChildren() {
	this.assertDisconnectResult(disconnected, disconnected)
	this.assertDisconnectResult(disconnecting, disconnecting) // don't interupt
	this.assertDisconnectResult(connected, disconnected)
	this.assertDisconnectResult(connecting, disconnected)
}
func (this *BrokerFixture) assertDisconnectResult(initial, updated uint64) {
	this.broker.state = initial
	this.broker.Disconnect()
	this.So(this.broker.State(), should.Equal, updated)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestDisconnectWithOnlyWriters() {
	this.broker.state = connected

	writers := []*FakeWriter{&FakeWriter{}, &FakeWriter{}}
	for _, writer := range writers {
		this.broker.writers = append(this.broker.writers, writer)
	}

	this.broker.Disconnect()
	this.broker.Disconnect() // only tries to shut down once

	this.So(writers[0].closed, should.Equal, 1)
	this.So(writers[1].closed, should.Equal, 1)
	this.So(this.broker.State(), should.Equal, disconnected)
	this.So(this.broker.writers, should.BeEmpty)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestDisconnectWithOnlyReaders() {
	this.broker.state = connected

	readers := []*FakeReader{&FakeReader{}, &FakeReader{}}
	for _, reader := range readers {
		this.broker.readers = append(this.broker.readers, reader)
	}

	this.broker.Disconnect()
	this.broker.Disconnect() // only tries to shut down once

	this.So(readers[0].closed, should.Equal, 1)
	this.So(readers[1].closed, should.Equal, 1)
	this.So(this.broker.State(), should.Equal, disconnecting)
	this.So(len(this.broker.readers), should.Equal, len(readers))
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestLastReaderShutdownComplete() {
	this.broker.state = disconnecting
	connection := &FakeConnection{}
	this.broker.connection = connection

	reader, writer := &FakeReader{}, &FakeWriter{}
	this.broker.readers = append(this.broker.readers, reader)
	this.broker.writers = append(this.broker.writers, writer)

	this.broker.removeReader(reader)

	this.So(this.broker.readers, should.BeEmpty)
	this.So(this.broker.writers, should.BeEmpty)
	this.So(writer.closed, should.Equal, 1)
	this.So(this.broker.State(), should.Equal, disconnected)
	this.So(this.broker.connection, should.BeNil)
	this.So(connection.closed, should.Equal, 1)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestSecondToLastReaderShutdownComplete() {
	this.broker.state = disconnecting

	reader1, reader2, writer := &FakeReader{}, &FakeReader{}, &FakeWriter{}
	this.broker.readers = append(this.broker.readers, reader1)
	this.broker.readers = append(this.broker.readers, reader2)
	this.broker.writers = append(this.broker.writers, writer)

	this.broker.removeReader(reader1)

	this.So(this.broker.readers, should.NotBeEmpty)
	this.So(this.broker.writers, should.NotBeEmpty)
	this.So(writer.closed, should.Equal, 0)
	this.So(this.broker.State(), should.Equal, disconnecting)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestIsolatedReaderCloseDoesntAffectBrokerState() {
	this.broker.state = connected
	reader := &FakeReader{}
	this.broker.readers = append(this.broker.readers, reader)

	this.broker.removeReader(reader)

	this.So(this.broker.readers, should.BeEmpty)
	this.So(this.broker.State(), should.Equal, connected)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestIsolatedWriterCloseDoesntAffectBrokerState() {
	this.broker.state = connected
	writer := &FakeWriter{}
	this.broker.writers = append(this.broker.writers, writer)

	this.broker.removeWriter(writer)

	this.So(this.broker.writers, should.BeEmpty)
	this.So(this.broker.State(), should.Equal, connected)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestOpenReaderDuringConnection() {
	this.assertValidReader(connecting)
	this.assertValidReader(connected)
}
func (this *BrokerFixture) assertValidReader(initialState uint64) {
	this.broker.state = initialState
	reader := this.broker.OpenReader("queue")
	this.So(reader, should.NotBeNil)
	this.So(reader.(*ChannelReader).queue, should.Equal, "queue")
	this.So(this.broker.readers, should.NotBeEmpty)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestOpenReaderDuringDisconnection() {
	this.assertNilReader(disconnecting)
	this.assertNilReader(disconnected)
}
func (this *BrokerFixture) assertNilReader(initialState uint64) {
	this.broker.state = initialState
	reader := this.broker.OpenReader("queue")
	this.So(reader, should.BeNil)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestOpenTransientReader() {
	this.broker.state = connecting
	bindings := []string{"1", "2"}

	reader := this.broker.OpenTransientReader(bindings)

	this.So(reader, should.NotBeNil)
	this.So(reader.(*ChannelReader).bindings, should.Resemble, bindings)
	this.So(this.broker.readers, should.NotBeEmpty)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestOpenWriterDuringConnection() {
	this.assertValidWriter(connecting)
	this.assertValidWriter(connected)
}
func (this *BrokerFixture) assertValidWriter(initialState uint64) {
	this.broker.state = initialState
	writer := this.broker.OpenWriter()
	this.So(writer, should.NotBeNil)
	this.So(reflect.TypeOf(writer).String(), should.Equal, "*rabbit.ChannelWriter")
	this.So(this.broker.writers, should.NotBeEmpty)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestOpenWriterDuringDisconnection() {
	this.assertNilWriter(disconnecting)
	this.assertNilWriter(disconnected)
}
func (this *BrokerFixture) assertNilWriter(initialState uint64) {
	this.broker.state = initialState
	writer := this.broker.OpenWriter()
	this.So(writer, should.BeNil)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestOpenTransactionalWriter() {
	this.broker.state = connecting

	writer := this.broker.OpenTransactionalWriter()

	this.So(writer, should.NotBeNil)
	this.So(reflect.TypeOf(writer).String(), should.Equal, "*rabbit.TransactionWriter")
	this.So(this.broker.writers, should.NotBeEmpty)
}

////////////////////////////////////////////////////////

func (this *BrokerFixture) TestOpenChannel() {
	this.broker.state = connecting

	channel := this.broker.openChannel()

	this.So(channel, should.NotBeNil)
	this.So(this.broker.state, should.Equal, connected)
}
func (this *BrokerFixture) TestNoChannelWhileDisconnected() {
	this.broker.state = disconnected
	this.So(this.broker.openChannel(), should.BeNil)

	this.broker.state = disconnecting
	this.So(this.broker.openChannel(), should.BeNil)
}
func (this *BrokerFixture) TestOpenChannelAfterUnderlyingConnectorFailureRetries() {
	this.connector = NewFakeConnector(1, 0)
	this.createBroker()
	this.broker.state = connecting

	channel := this.broker.openChannel()

	this.So(channel, should.NotBeNil)
	this.So(this.connector.attempts, should.BeGreaterThan, 1)
	this.So(this.broker.state, should.Equal, connected)
	this.So(this.sleeper.Naps[0], should.Equal, time.Second*4)
}
func (this *BrokerFixture) TestOpenChannelAfterUnderlyingConnectionFailureRetries() {
	this.connector = NewFakeConnector(0, 1)
	this.createBroker()
	this.broker.state = connecting

	channel := this.broker.openChannel()

	this.So(channel, should.NotBeNil)
	this.So(this.connector.attempts, should.Equal, 2)
	this.So(this.connector.connection.attempts, should.Equal, 2)
	this.So(this.broker.state, should.Equal, connected)
	this.So(this.broker.connection, should.NotBeNil)
	this.So(this.sleeper.Naps[0], should.Equal, time.Second*4)
}
func (this *BrokerFixture) TestOpenChannelClosesConnectionOnFailure() {
	this.connector = NewFakeConnector(0, 2)
	this.createBroker()
	this.broker.state = connecting

	channel := this.broker.openChannel()

	this.So(channel, should.NotBeNil)
	this.So(this.connector.attempts, should.Equal, 3)
	this.So(this.connector.connection.attempts, should.Equal, 3)
	this.So(this.connector.connection.closed, should.Equal, 2)
	this.So(this.broker.state, should.Equal, connected)
	this.So(this.broker.connection, should.NotBeNil)
	this.So(this.sleeper.Naps[0], should.Equal, time.Second*4)
	this.So(this.sleeper.Naps[1], should.Equal, time.Second*4)
}

////////////////////////////////////////////////////////

type FakeConnector struct {
	attempts   int
	target     url.URL
	failures   int
	connection *FakeConnection
}

func NewFakeConnector(connectorFailures, connectionFailures int) *FakeConnector {
	return &FakeConnector{
		failures:   connectorFailures,
		connection: &FakeConnection{failures: connectionFailures},
	}
}

func (this *FakeConnector) Connect(target url.URL) (Connection, error) {
	this.attempts++
	if this.failures >= this.attempts {
		return nil, errors.New("Fail!")
	}

	return this.connection, nil
}

////////////////////////////////////////////////////////

type FakeConnection struct {
	attempts int
	failures int
	closed   int
}

func (this *FakeConnection) Channel() (Channel, error) {
	this.attempts++
	if this.failures >= this.attempts {
		return nil, errors.New("Fail")
	}

	return &FakeChannel{}, nil
}

func (this *FakeConnection) Close() error {
	this.closed++
	return nil
}

////////////////////////////////////////////////////////

type FakeChannel struct{}

func (this *FakeChannel) ConfigureChannelBuffer(int) error                     { return nil }
func (this *FakeChannel) DeclareTransientQueue() (string, error)               { return "", nil }
func (this *FakeChannel) BindExchangeToQueue(string, string) error             { return nil }
func (this *FakeChannel) Consume(string, string) (<-chan amqp.Delivery, error) { return nil, nil }
func (this *FakeChannel) ExclusiveConsume(string, string) (<-chan amqp.Delivery, error) {
	return nil, nil
}
func (this *FakeChannel) CancelConsumer(string) error                  { return nil }
func (this *FakeChannel) Close() error                                 { return nil }
func (this *FakeChannel) AcknowledgeSingleMessage(uint64) error        { return nil }
func (this *FakeChannel) AcknowledgeMultipleMessages(uint64) error     { return nil }
func (this *FakeChannel) ConfigureChannelAsTransactional() error       { return nil }
func (this *FakeChannel) PublishMessage(string, amqp.Publishing) error { return nil }
func (this *FakeChannel) CommitTransaction() error                     { return nil }
func (this *FakeChannel) RollbackTransaction() error                   { return nil }

////////////////////////////////////////////////////////

type FakeWriter struct{ closed int }

func (this *FakeWriter) Close()                         { this.closed++ }
func (this *FakeWriter) Write(messenger.Dispatch) error { return nil }

////////////////////////////////////////////////////////

type FakeReader struct{ closed int }

func (this *FakeReader) Close()                                { this.closed++ }
func (this *FakeReader) Listen()                               {}
func (this *FakeReader) Deliveries() <-chan messenger.Delivery { return nil }
func (this *FakeReader) Acknowledgements() chan<- interface{}  { return nil }
