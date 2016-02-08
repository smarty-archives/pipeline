package transports

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"time"
)

type ChannelSelectWriterFixture struct {
	*gunit.Fixture

	writer     *ChannelSelectWriter
	fakeWriter *FakeWriter
}

func (this *ChannelSelectWriterFixture) Setup() {
	this.fakeWriter = &FakeWriter{}
	this.writer = NewChannelSelectWriter(this.fakeWriter, 2)
}

func (this *ChannelSelectWriterFixture) TestWritesReachInnerWriter() {
	go this.writer.Listen()

	written1, err1 := this.writer.Write([]byte("Hello,"))
	written2, err2 := this.writer.Write([]byte("World!"))
	time.Sleep(time.Millisecond)

	this.So(written1, should.Equal, len("Hello,"))
	this.So(written2, should.Equal, len("World!"))

	this.So(err1, should.BeNil)
	this.So(err2, should.BeNil)

	this.So(this.fakeWriter.Writes(), should.Equal, 2)
	this.So(string(this.fakeWriter.written[0]), should.Equal, "Hello,")
	this.So(string(this.fakeWriter.written[1]), should.Equal, "World!")
}

func (this *ChannelSelectWriterFixture) TestEmptyBuffersNeverReachInnerWriter() {
	go this.writer.Listen()

	written, err := this.writer.Write([]byte(""))
	time.Sleep(time.Millisecond)

	this.So(written, should.Equal, 0)
	this.So(err, should.BeNil)
	this.So(this.fakeWriter.Writes(), should.Equal, 0)
}

func (this *ChannelSelectWriterFixture) TestWritesWhichOverflowChannelAreDiscarded() {
	this.writer.Write([]byte("a"))
	this.writer.Write([]byte("b"))
	written, err := this.writer.Write([]byte("c"))

	this.So(written, should.Equal, 0)
	this.So(err, should.Equal, WriteDiscardedError)
}

//////////////////////////////////////////////////////////////////////////

type FakeWriter struct {
	written [][]byte
}

func (this *FakeWriter) Writes() int { return len(this.written) }
func (this *FakeWriter) Write(buffer []byte) (int, error) {
	this.written = append(this.written, buffer)
	return 0, nil
}
