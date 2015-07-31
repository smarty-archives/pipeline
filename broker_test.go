package rabbit

import (
	"errors"
	"net/url"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type BrokerFixture struct {
	*gunit.Fixture

	target    url.URL
	connector *FakeConnector
	broker    *Broker
}

func (this *BrokerFixture) Setup() {
	target, _ := url.Parse("amqps://username:password@localhost:5671/")
	this.target = *target
	this.createBroker()
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
	this.So(connection.closed, should.BeTrue)
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

type FakeConnector struct {
	attempts           int
	target             url.URL
	connectorFailures  int
	connectionFailures int
}

func NewFakeConnector(connectorFailures, connectionFailures int) *FakeConnector {
	return &FakeConnector{
		connectorFailures:  connectorFailures,
		connectionFailures: connectionFailures,
	}
}

func (this *FakeConnector) Connect(target url.URL) (Connection, error) {
	this.attempts++
	if this.connectorFailures > this.attempts {
		return nil, errors.New("Fail!")
	}

	return NewFakeConnection(), nil
}

////////////////////////////////////////////////////////

type FakeConnection struct {
	attempts int
	failures int
	closed   bool
}

func NewFakeConnection() *FakeConnection {
	return &FakeConnection{}
}

func (this *FakeConnection) Channel() (Channel, error) {
	this.attempts++
	if this.failures > this.attempts {
		return nil, errors.New("Fail")
	}

	return nil, nil
}

func (this *FakeConnection) Close() error {
	this.closed = true
	return nil
}

////////////////////////////////////////////////////////

type FakeWriter struct{ closed int }

func (this *FakeWriter) Close()               { this.closed++ }
func (this *FakeWriter) Write(Dispatch) error { return nil }

////////////////////////////////////////////////////////

type FakeReader struct{ closed int }

func (this *FakeReader) Close()                               { this.closed++ }
func (this *FakeReader) Listen()                              {}
func (this *FakeReader) Deliveries() <-chan Delivery          { return nil }
func (this *FakeReader) Acknowledgements() chan<- interface{} { return nil }
