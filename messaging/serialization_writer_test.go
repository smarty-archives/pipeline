package messaging

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type SerializationWriterFixture struct {
	*gunit.Fixture

	writer     *SerializationWriter
	inner      *FakeCommitWriter
	serializer *FakeSerializer
	discovery  *FakeDiscovery
}

func (this *SerializationWriterFixture) Setup() {
	log.SetOutput(ioutil.Discard)

	this.inner = &FakeCommitWriter{}
	this.serializer = &FakeSerializer{}
	this.discovery = &FakeDiscovery{}
	this.buildWriter()
}
func (this *SerializationWriterFixture) Teardown() {
	log.SetOutput(os.Stdout)
}

func (this *SerializationWriterFixture) buildWriter() {
	this.writer = NewSerializationWriter(this.inner, this.serializer, this.discovery)
}

////////////////////////////////////////////////////////////////////////////////

func (this *SerializationWriterFixture) TestWriterAddsSerializedPayloadAndTypeToDispatch() {
	this.inner.writeError = errors.New("ensure inner errors are returned to caller")

	original := Dispatch{
		SourceID:    1,
		MessageID:   2,
		Destination: "3",
		MessageType: "", // should be populated
		Encoding:    "4",
		Durable:     true,
		Expiration:  time.Now(),
		Payload:     nil,
		Message:     TestMessage{},
	}

	err := this.writer.Write(original)

	this.So(err, should.Equal, this.inner.writeError)
	this.So(this.inner.written, should.NotBeEmpty)
	this.So(this.inner.written[0], should.Resemble, Dispatch{
		SourceID:    1,
		MessageID:   2,
		Destination: "3",
		MessageType: "message.type.discovered",
		Encoding:    "4",
		Durable:     true,
		Expiration:  original.Expiration,
		Payload:     testSerializedPayload,
		Message:     TestMessage{},
	})
}

////////////////////////////////////////////////////////////////////////////////

func (this *SerializationWriterFixture) TestSerializationFails() {
	this.serializer.serializeError = errors.New("Serialization failed")
	this.buildWriter()

	err := this.writer.Write(Dispatch{Message: TestMessage{}})
	this.So(err, should.Equal, this.serializer.serializeError)
	this.So(this.inner.written, should.BeEmpty)
}

func (this *SerializationWriterFixture) TestDispatchAlreadyContainsSerializedPayload() {
	this.inner.writeError = errors.New("ensure inner errors are returned to caller")

	message := Dispatch{
		MessageType: "untouched",
		Payload:     []byte("already serialized"),
	}
	err := this.writer.Write(message)
	this.So(err, should.Equal, this.inner.writeError)
	this.So(this.serializer.called, should.Equal, 0)
	this.So(this.inner.written, should.Resemble, []Dispatch{message})
}

////////////////////////////////////////////////////////////////////////////////

func (this *SerializationWriterFixture) TestDispatchAlreadyContainsMessageType() {
	message := Dispatch{
		MessageType: "untouched",
		Message:     TestMessage{},
	}
	this.writer.Write(message)
	this.So(this.inner.written, should.Resemble, []Dispatch{Dispatch{
		MessageType: "untouched",
		Payload:     testSerializedPayload,
		Message:     TestMessage{},
	}})
}

////////////////////////////////////////////////////////////////////////////////

func (this *SerializationWriterFixture) TestMessageTypeDiscoveryErrorsReturned() {
	this.discovery.discoveryError = errors.New("discovery error")
	message := Dispatch{Message: TestMessage{}}
	err := this.writer.Write(message)
	this.So(err, should.Equal, this.discovery.discoveryError)
	this.So(this.inner.written, should.BeEmpty)
}

////////////////////////////////////////////////////////////////////////////////

func (this *SerializationWriterFixture) TestWriteEmptyDispatch() {
	err := this.writer.Write(Dispatch{})
	this.So(err, should.Equal, EmptyDispatchError)
}

////////////////////////////////////////////////////////////////////////////////

func (this *SerializationWriterFixture) TestWriteEmptyMessageType() {
	err := this.writer.Write(Dispatch{Payload: []byte("already serialized")})
	this.So(err, should.Equal, MessageTypeDiscoveryError)
	this.So(this.inner.written, should.BeEmpty)
}

////////////////////////////////////////////////////////////////////////////////

func (this *SerializationWriterFixture) TestCommitInvokesUnderlyingWriter() {
	this.inner.commitError = errors.New("commit error")
	err := this.writer.Commit()
	this.So(this.inner.commits, should.Equal, 1)
	this.So(err, should.Equal, this.inner.commitError)
}

////////////////////////////////////////////////////////////////////////////////

func (this *SerializationWriterFixture) TestCommitOnRegularWriterPanics() {
	this.writer = NewSerializationWriter(&FakeWriter{}, this.serializer, this.discovery)
	err := this.writer.Commit()
	this.So(err, should.BeNil)
}

////////////////////////////////////////////////////////////////////////////////

func (this *SerializationWriterFixture) TestCloseUnderlyingWriter() {
	this.writer.Close()
	this.So(this.inner.closed, should.Equal, 1)
}

////////////////////////////////////////////////////////////////////////////////

type FakeWriter struct{}

func (this *FakeWriter) Write(dispatch Dispatch) error { return nil }
func (this *FakeWriter) Close()                        {}

////////////////////////////////////////////////////////////////////////////////

type FakeCommitWriter struct {
	written     []Dispatch
	commits     int
	closed      int
	writeError  error
	commitError error
}

func (this *FakeCommitWriter) Write(dispatch Dispatch) error {
	this.written = append(this.written, dispatch)
	return this.writeError
}

func (this *FakeCommitWriter) Commit() error {
	this.commits++
	return this.commitError
}

func (this *FakeCommitWriter) Close() {
	this.closed++
}

////////////////////////////////////////////////////////////////////////////////

type FakeSerializer struct {
	called         int
	serializeError error
}

func (this *FakeSerializer) Serialize(interface{}) ([]byte, error) {
	this.called++
	if this.serializeError != nil {
		return nil, this.serializeError
	} else {
		return testSerializedPayload, nil
	}
}

var testSerializedPayload = []byte("serializer called successfully")

////////////////////////////////////////////////////////////////////////////////

type FakeDiscovery struct {
	discoveryError error
}

func (this *FakeDiscovery) Discover(instance interface{}) (string, error) {
	if this.discoveryError == nil {
		return "message.type.discovered", nil
	} else {
		return "", this.discoveryError
	}
}

////////////////////////////////////////////////////////////////////////////////

type TestMessage struct{}

////////////////////////////////////////////////////////////////////////////////
