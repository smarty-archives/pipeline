package pipeline

import (
	"reflect"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/go-messenger"
	"github.com/smartystreets/gunit"
)

type JSONDeserializerFixture struct {
	*gunit.Fixture

	deserializer *JSONDeserializer
}

func (this *JSONDeserializerFixture) Setup() {
	this.deserializer = NewJSONDeserializer(map[string]reflect.Type{
		"ApplicationEvent": reflect.TypeOf(ApplicationEvent{}),
	})
}

func (this *JSONDeserializerFixture) TestDeserializeKnownMessageType() {
	delivery := &messenger.Delivery{
		MessageType: "ApplicationEvent",
		Payload:     []byte(`{"Message": "Hello, World!"}`),
	}

	this.deserializer.Deserialize(delivery)

	this.So(delivery.Message, should.Resemble, &ApplicationEvent{Message: "Hello, World!"})
}

func (this *JSONDeserializerFixture) TestDeserializeUnknownMessageType() {
	delivery := &messenger.Delivery{
		MessageType: "What Am I?",
		Payload:     []byte(`{"Message": "Hello, World!"}`),
	}

	this.deserializer.Deserialize(delivery)

	this.So(delivery.Message, should.BeNil)
}

func (this *JSONDeserializerFixture) TestDeserializeUnkonwnMessageTypeIsCriticalFailure() {
	this.deserializer.PanicWhenMessageTypeIsUnknown()

	delivery := &messenger.Delivery{
		MessageType: "What Am I?",
		Payload:     []byte(`{"Message": "Hello, World!"}`),
	}

	this.So(func() { this.deserializer.Deserialize(delivery) }, should.Panic)
}

func (this *JSONDeserializerFixture) TestDeserializationFailsSilently() {
	delivery := &messenger.Delivery{
		MessageType: "ApplicationEvent",
		Payload:     []byte(`ThisIsNotJSON`),
	}

	this.deserializer.Deserialize(delivery)

	this.So(delivery.Message, should.BeNil)
}

func (this *JSONDeserializerFixture) TestDeserializationFailsWithPanic() {
	this.deserializer.PanicWhenDeserializationFails()

	delivery := &messenger.Delivery{
		MessageType: "ApplicationEvent",
		Payload:     []byte(`blah blah blah`),
	}

	this.So(func() { this.deserializer.Deserialize(delivery) }, should.Panic)
}

///////////////////////////////////////////////////////////////

type ApplicationEvent struct {
	Message string
}
