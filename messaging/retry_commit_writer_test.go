package messaging

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestRetryCommitWriterFixture(t *testing.T) {
	gunit.Run(new(RetryCommitWriterFixture), t)
}

type RetryCommitWriterFixture struct {
	*gunit.Fixture

	writer *RetryCommitWriter
	inner  *FakeRetryCommitWriter

	successes    int
	successInput uint64
	sleeps       int
	sleepInput   uint64
}

func (this *RetryCommitWriterFixture) Setup() {
	this.inner = &FakeRetryCommitWriter{}
	this.writer = NewRetryCommitWriter(this.inner, 42, this.success, this.sleep)
}
func (this *RetryCommitWriterFixture) success(value uint64) {
	this.successes++
	this.successInput = value
}
func (this *RetryCommitWriterFixture) sleep(value uint64) {
	this.sleeps++
	this.sleepInput = value
}

///////////////////////////////////////////////////////////////

func (this *RetryCommitWriterFixture) TestNoErrorsNoRetries() {
	this.inner.errorsUntil = 0
	dispatches := []Dispatch{Dispatch{}, Dispatch{}, Dispatch{}}

	for _, item := range dispatches {
		this.writer.Write(item)
	}

	err := this.writer.Commit()

	this.So(err, should.BeNil)
	this.So(this.inner.written, should.Resemble, dispatches)
	this.So(this.inner.writes, should.Equal, len(dispatches))
	this.So(this.inner.commits, should.Equal, 1)
	this.So(this.sleeps, should.Equal, 0)
	this.So(this.sleepInput, should.Equal, 0)
	this.So(this.successes, should.Equal, 1)
	this.So(this.successInput, should.Equal, 0)
}

///////////////////////////////////////////////////////////////

func (this *RetryCommitWriterFixture) TestRetryUntilNoErrors() {
	this.inner.errorsUntil = 2
	dispatches := []Dispatch{Dispatch{}, Dispatch{}, Dispatch{}}

	for _, item := range dispatches {
		this.writer.Write(item)
	}

	err := this.writer.Commit()

	this.So(err, should.BeNil)
	this.So(this.inner.written[0:3], should.Resemble, dispatches)
	this.So(this.inner.written[3:6], should.Resemble, dispatches)
	this.So(this.inner.written[6:9], should.Resemble, dispatches)
	this.So(this.inner.writes, should.Equal, len(dispatches)*3)
	this.So(this.inner.commits, should.Equal, 3)
	this.So(this.sleeps, should.Equal, 2)
	this.So(this.sleepInput, should.Equal, 1)
	this.So(this.writer.buffer, should.BeEmpty)
	this.So(this.successes, should.Equal, 1)
	this.So(this.successInput, should.Equal, 2)
}

///////////////////////////////////////////////////////////////

func (this *RetryCommitWriterFixture) TestRetryUntilClosed() {
	dispatches := []Dispatch{Dispatch{}, Dispatch{}, Dispatch{}}
	this.writer.Close()

	for _, item := range dispatches {
		this.writer.Write(item)
	}

	err := this.writer.Commit()

	this.So(err, should.Equal, WriterClosedError)
	this.So(this.inner.written, should.Resemble, dispatches[0:1])
	this.So(this.inner.writes, should.Equal, 1)
	this.So(this.inner.commits, should.Equal, 0)
	this.So(this.sleeps, should.Equal, 0)
	this.So(this.sleepInput, should.Equal, 0)
	this.So(this.successes, should.Equal, 0)
	this.So(this.successInput, should.Equal, 0)
}

///////////////////////////////////////////////////////////////

type FakeRetryCommitWriter struct {
	writes      int
	closed      int
	commits     int
	errorsUntil int
	written     []Dispatch
}

func (this *FakeRetryCommitWriter) Write(message Dispatch) error {
	this.writes++
	this.written = append(this.written, message)

	if this.closed > 0 {
		return WriterClosedError
	} else if this.errorsUntil >= this.commits {
		return errors.New("general write failure")
	}

	return nil
}
func (this *FakeRetryCommitWriter) Commit() error {
	this.commits++

	if this.errorsUntil >= this.commits {
		return errors.New("Unable to commit")
	}

	return nil
}

func (this *FakeRetryCommitWriter) Close() {
	this.closed++
}
