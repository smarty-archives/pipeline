package rabbit

import (
	"errors"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/pipeline/messaging"
)

type TransactionWriterFixture struct {
	*gunit.Fixture

	writer     *TransactionWriter
	controller *FakeWriterController
}

func (this *TransactionWriterFixture) Setup() {
	this.controller = newFakeWriterController()
	this.writer = transactionWriter(this.controller)
}

///////////////////////////////////////////////////////////////

func (this *TransactionWriterFixture) TestDispatchIsWrittenToChannel() {
	dispatch := messaging.Dispatch{
		Destination: "destination",
		Payload:     []byte{1, 2, 3, 4, 5},
	}

	this.writer.Write(dispatch)

	this.So(this.controller.channel.exchange, should.Equal, dispatch.Destination)
	this.So(this.controller.channel.dispatch.Body, should.Resemble, dispatch.Payload)
	this.So(this.controller.channel.transactional, should.BeTrue)
}

///////////////////////////////////////////////////////////////

func (this *TransactionWriterFixture) TestChannelCannotBeObtained() {
	this.controller.channel = nil

	this.So(func() { this.writer.Write(messaging.Dispatch{}) }, should.NotPanic)
}

///////////////////////////////////////////////////////////////

func (this *TransactionWriterFixture) TestFailedChannelNOTClosedOnFailedWrites() {
	this.controller.channel.err = errors.New("channel failed")

	this.writer.Write(messaging.Dispatch{})

	this.So(this.controller.channel.closed, should.Equal, 0)
	this.So(this.writer.channel, should.NotBeNil)
}

///////////////////////////////////////////////////////////////

func (this *TransactionWriterFixture) TestCloseWriter() {
	this.writer.Close()
	this.So(this.writer.closed, should.BeTrue)
}

///////////////////////////////////////////////////////////////

func (this *TransactionWriterFixture) TestCommitWithoutWriteFails() {
	err := this.writer.Commit()

	this.So(err, should.Equal, commitBeforeWriteError)
}

func (this *TransactionWriterFixture) TestCommitCallsUnderlyingChannel() {
	this.writer.Write(messaging.Dispatch{})

	err := this.writer.Commit()
	this.So(err, should.BeNil)
	this.So(this.controller.channel.commits, should.Equal, 1)
}

func (this *TransactionWriterFixture) TestFailedCommitsReturnError() {
	this.writer.Write(messaging.Dispatch{})
	this.controller.channel.err = errors.New("Commit failure")

	err := this.writer.Commit()
	this.So(err, should.Equal, this.controller.channel.err)
	this.So(this.controller.channel.commits, should.Equal, 1)
	this.So(this.controller.channel.closed, should.Equal, 1)
	this.So(this.writer.channel, should.BeNil)
}
