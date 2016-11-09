package messaging

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"
)

func TestJSONSerializerFixture(t *testing.T) {
	gunit.Run(new(JSONSerializerFixture), t)
}

type JSONSerializerFixture struct {
	*gunit.Fixture

	serializer *JSONSerializer
}

func (this *JSONSerializerFixture) Setup() {
	log.SetOutput(ioutil.Discard)
	this.serializer = NewJSONSerializer()
	this.serializer.logger = logging.Capture()
}

func (this *JSONSerializerFixture) TestSerializationSucceeds() {
	message := ExampleMessage{Content: "Hello, World!"}
	content, err := this.serializer.Serialize(message)
	this.So(err, should.BeNil)
	this.So(string(content), should.Equal, `{"Content":"Hello, World!"}`)
}

func (this *JSONSerializerFixture) TestSerializationFails() {
	message := InvalidMessage{Stuff: make(chan int, 4)}
	content, err := this.serializer.Serialize(message)
	this.So(err, should.NotBeNil)
	this.So(content, should.BeNil)
}

func (this *JSONSerializerFixture) TestContentType() {
	this.So(this.serializer.ContentType(), should.Equal, "application/json")
}

func (this *JSONSerializerFixture) TestContentEncoding() {
	this.So(this.serializer.ContentEncoding(), should.Equal, "")
}

func (this *JSONSerializerFixture) TestSerializationFailsAndPanics() {
	this.serializer.PanicWhenSerializationFails()
	message := InvalidMessage{Stuff: make(chan int, 4)}
	this.So(func() { this.serializer.Serialize(message) }, should.Panic)
}

////////////////////////////////////////////////////////////////////////////////

type ExampleMessage struct {
	Content string
}

type InvalidMessage struct {
	Stuff chan int `json:"stuff"`
}
