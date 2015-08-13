package messaging

import (
	"io/ioutil"
	"log"

	"github.com/smartystreets/gunit"
)

type SerializationWriterFixture struct {
	*gunit.Fixture

	writer     *SerializationWriter
	serializer *FakeSerializer
}

func (this *SerializationWriterFixture) Setup() {
	log.SetOutput(ioutil.Discard)
	this.serializer = &FakeSerializer{}
	this.writer = NewSerializationWriter(this.serializer)
}

// TODO: does this writer receive a Dispatch or a Delivery? If, like all other
// writers it receives a Dispatch, how does it get the actual Message? (which
// is on the Delivery...)

////////////////////////////////////////////////////////////////////////////////

type FakeSerializer struct {
	called int
}

func (this *FakeSerializer) Serialize(interface{}) ([]byte, error) {
	this.called++
	return []byte("X"), nil
}
