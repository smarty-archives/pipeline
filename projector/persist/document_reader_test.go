package persist

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestDocumentReaderFixture(t *testing.T) {
	gunit.Run(new(DocumentReaderFixture), t)
}

type DocumentReaderFixture struct {
	*gunit.Fixture

	path     string
	reader   *DocumentReader
	client   *FakeHTTPGetClient // HTTPClient
	document *Document
}

func (this *DocumentReaderFixture) Setup() {
	this.path = "/document/path"
	this.client = &FakeHTTPGetClient{}
	this.reader = NewDocumentReader(this.client)
	this.document = &Document{}
}

func (this *DocumentReaderFixture) TestRequestInvalid_ClientIgnored() {
	this.path = "%%%%%%%%"
	this.assertPanic(`Could not create request: parse %%%%%%%%: invalid URL escape "%%%"`)
	this.So(this.client.called, should.BeFalse)
}

func (this *DocumentReaderFixture) TestClientErrorPreventsDocumentReading() {
	this.client.err = errors.New("BOINK!")
	this.assertPanic("HTTP Client Error: BOINK!")
}

func (this *DocumentReaderFixture) TestDocumentNotFound_JSONMarshalNotAttempted() {
	this.client.response = &http.Response{StatusCode: 404, Body: newHTTPBody("Not found")}
	this.read()
	this.So(this.document.ID, should.Equal, 0)
}

func (this *DocumentReaderFixture) TestBodyUnreadable() {
	var BodyUnreadableResponse = &http.Response{StatusCode: 200, Body: newReadErrorHTTPBody()}
	this.client.response = BodyUnreadableResponse
	this.So(this.read, should.Panic)
	this.So(this.document.ID, should.Equal, 0)
	this.So(BodyUnreadableResponse.Body.(*FakeHTTPResponseBody).closed, should.BeTrue)
}

func (this *DocumentReaderFixture) TestBadJSON() {
	var BadJSONResponse = &http.Response{StatusCode: 200, Body: newHTTPBody("I am bad JSON.")}
	this.client.response = BadJSONResponse
	this.So(this.read, should.Panic)
	this.So(this.document.ID, should.Equal, 0)
	this.So(BadJSONResponse.Body.(*FakeHTTPResponseBody).closed, should.BeTrue)
}

func (this *DocumentReaderFixture) TestValidUncompressedResponse_PopulatesDocument() {
	var ValidUncompressedResponse = &http.Response{StatusCode: 200, Body: newHTTPBody(`{"ID": 1234}`)}
	this.client.response = ValidUncompressedResponse
	this.read()
	this.So(this.document.ID, should.Equal, 1234)
	this.So(ValidUncompressedResponse.Body.(*FakeHTTPResponseBody).closed, should.BeTrue)
}
func (this *DocumentReaderFixture) TestValidCompressedResponse_PopulatesDocument() {
	var ValidCompressedResponse = &http.Response{StatusCode: 200, Body: newHTTPBody(`{"ID": 1234}`)}

	ValidCompressedResponse.Header = make(http.Header)
	ValidCompressedResponse.Header.Set("Content-Encoding", "gzip")

	targetBuffer := bytes.NewBuffer([]byte{})
	writer := gzip.NewWriter(targetBuffer)
	io.Copy(writer, ValidCompressedResponse.Body)
	writer.Close()

	ValidCompressedResponse.Body = ioutil.NopCloser(targetBuffer)

	this.client.response = ValidCompressedResponse
	this.read()
	this.So(this.document.ID, should.Equal, 1234)
}
func (this *DocumentReaderFixture) read() {
	this.reader.ReadPanic(this.path, this.document)
}
func (this *DocumentReaderFixture) assertPanic(message string) {
	this.So(this.read, should.Panic)
	this.So(this.document.ID, should.Equal, 0)
}

//////////////////////////////////////////////////////////////////////////////////////////////

type FakeHTTPGetClient struct {
	err      error
	response *http.Response
	called   bool
}

func (this *FakeHTTPGetClient) Do(request *http.Request) (*http.Response, error) {
	this.called = true
	return this.response, this.err
}

///////////////////////////////////////////////////////////////////////////////////////////

type Document struct{ ID int }

////////////////////////////////////////////////////////////////////////////////////////////

func newHTTPBody(message string) io.ReadCloser {
	return &FakeHTTPResponseBody{Reader: strings.NewReader(message)}
}
func newReadErrorHTTPBody() io.ReadCloser {
	return &FakeHTTPResponseBody{err: errors.New("BOINK!")}
}

type FakeHTTPResponseBody struct {
	*strings.Reader

	err    error
	closed bool
}

func (this *FakeHTTPResponseBody) Read(p []byte) (int, error) {
	if this.err != nil {
		return 0, this.err
	}
	return this.Reader.Read(p)
}

func (this *FakeHTTPResponseBody) Close() error {
	this.closed = true
	return nil
}
