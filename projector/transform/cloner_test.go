package transform

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"
	"github.com/smartystreets/pipeline/projector"
)

func TestClonerFixture(t *testing.T) {
	gunit.Run(new(ClonerFixture), t)
}

type ClonerFixture struct {
	*gunit.Fixture

	buffer   ResetReadWriter
	cloner   *DocumentCloner
	original projector.Document
	clone    projector.Document
}

func (this *ClonerFixture) Setup() {
	this.buffer = &bytes.Buffer{}
	this.initializeDocumentCloner()
	this.original = NewCloneableReport(42)
}

func (this *ClonerFixture) initializeDocumentCloner() {
	this.cloner = NewDocumentCloner(this.buffer)
	this.cloner.logger = logging.Capture()
}

func (this *ClonerFixture) Clone() {
	this.clone = this.cloner.Clone(this.original)
}

func (this *ClonerFixture) FillBufferWithGarbage() {
	for x := 0; x < 1024; x++ {
		fmt.Fprint(this.buffer, x)
	}
}

//////////////////////////////////////////////////////////////////

func (this *ClonerFixture) TestClonedDocumentResembledOriginalButIsNOTTheOriginal() {
	this.Clone()
	this.So(this.clone, should.Resemble, this.original)
	this.So(this.clone, should.NotPointTo, this.original)
}

//////////////////////////////////////////////////////////////////

func (this *ClonerFixture) TestClonePanicsIfGOBEncodingFails() {
	this.original = NewUncloneableReport(42)
	this.So(this.Clone, should.Panic)
}

//////////////////////////////////////////////////////////////////

func (this *ClonerFixture) TestClonePanicsIfGOBDecodingFails() {
	this.buffer = &EOFReadBuffer{Buffer: &bytes.Buffer{}}
	this.initializeDocumentCloner()
	this.So(this.Clone, should.Panic)
}

//////////////////////////////////////////////////////////////////

func (this *ClonerFixture) TestCloneShouldOperateOnAnEmptyBufferEachTime() {
	for x := 0; x < 10; x++ {
		this.FillBufferWithGarbage()
		this.So(this.Clone, should.NotPanic)
	}
}

//////////////////////////////////////////////////////////////////

type CloneableReport struct{ ID int }

func NewCloneableReport(id int) *CloneableReport                            { return &CloneableReport{ID: id} }
func (this *CloneableReport) Path() string                                  { panic("NOT IMPLEMENTED") }
func (this *CloneableReport) Apply(message interface{}) bool                { panic("NOT IMPLEMENTED") }
func (this *CloneableReport) Lapse(now time.Time) (next projector.Document) { panic("NOT IMPLEMENTED") }

//////////////////////////////////////////////////////////////////

// This type has no exported fields, which means it cannot be encoded by the gob.Encoder.
type UncloneableReport struct{ id int }

func NewUncloneableReport(id int) *UncloneableReport           { return &UncloneableReport{id: id} }
func (this *UncloneableReport) Path() string                   { panic("NOT IMPLEMENTED") }
func (this *UncloneableReport) Apply(message interface{}) bool { panic("NOT IMPLEMENTED") }
func (this *UncloneableReport) Lapse(now time.Time) (next projector.Document) {
	panic("NOT IMPLEMENTED")
}

///////////////////////////////////////////////////////////////////

type EOFReadBuffer struct{ *bytes.Buffer }

// Read only returns io.EOF so that we can expose an error condition in the gob.Decoder.
func (this *EOFReadBuffer) Read(p []byte) (n int, err error) { return 0, io.EOF }

///////////////////////////////////////////////////////////////////
