package persist

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type GetRetryClientFixture struct {
	*gunit.Fixture

	fakeClient  *FakeHTTPClientForGetRetry
	retryClient *GetRetryClient
	response    *http.Response
	err         error
	naps        int
	stdout      *bytes.Buffer
}

func (this *GetRetryClientFixture) Setup() {
	this.stdout = new(bytes.Buffer)
	log.SetOutput(this.stdout)
	napTime = func(time.Duration) { this.naps++ }
	this.fakeClient = &FakeHTTPClientForGetRetry{}
	this.retryClient = NewGetRetryClient(this.fakeClient, retries)
}

/////////////////////////////////////////////////////////

func (this *GetRetryClientFixture) TestClientFindsDocumentOnFirstTry() {
	this.fakeClient.statusCode = http.StatusOK
	request, _ := http.NewRequest("GET", "/document", nil)
	this.response, this.err = this.retryClient.Do(request)
	if this.So(this.response, should.NotBeNil) {
		this.So(this.response.StatusCode, should.Equal, http.StatusOK)
	}
	this.So(this.err, should.BeNil)
}

/////////////////////////////////////////////////////////

func (this *GetRetryClientFixture) TestClientFindsNODocumentOnFirstTry() {
	this.fakeClient.statusCode = http.StatusNotFound
	request, _ := http.NewRequest("GET", "/document", nil)
	this.response, this.err = this.retryClient.Do(request)
	if this.So(this.response, should.NotBeNil) {
		this.So(this.response.StatusCode, should.Equal, http.StatusNotFound)
	}
	this.So(this.err, should.BeNil)
}

/////////////////////////////////////////////////////////

func (this *GetRetryClientFixture) TestClientFailsAtFirst_ThenSucceeds() {
	this.fakeClient.statusCode = http.StatusOK
	request, _ := http.NewRequest("GET", "/fail-first", nil)
	this.response, this.err = this.retryClient.Do(request)
	if this.So(this.response, should.NotBeNil) {
		this.So(this.response.StatusCode, should.Equal, http.StatusOK)
	}
	this.So(this.err, should.BeNil)
}

/////////////////////////////////////////////////////////

func (this *GetRetryClientFixture) TestClientFailsAtFirst_ThenFindsNoDocument() {
	this.fakeClient.statusCode = http.StatusNotFound
	request, _ := http.NewRequest("GET", "/fail-first", nil)
	this.response, this.err = this.retryClient.Do(request)
	if this.So(this.response, should.NotBeNil) {
		this.So(this.response.StatusCode, should.Equal, http.StatusNotFound)
	}
	this.So(this.err, should.BeNil)
}

/////////////////////////////////////////////////////////

func (this *GetRetryClientFixture) TestClientNeverSucceeds() {
	request, _ := http.NewRequest("GET", "/fail-always", nil)
	this.response, this.err = this.retryClient.Do(request)
	this.So(this.response, should.BeNil)
	this.So(this.err, should.NotBeNil)
	this.So(this.fakeClient.calls, should.Equal, maxAttempts)
	this.So(this.naps, should.Equal, maxAttempts)
}

/////////////////////////////////////////////////////////

func (this *GetRetryClientFixture) TestClientBadStatusCodeAtFirst_ThenFindsDocument() {
	this.fakeClient.statusCode = http.StatusOK
	request, _ := http.NewRequest("GET", "/bad-status", nil)
	this.response, this.err = this.retryClient.Do(request)
	if this.So(this.response, should.NotBeNil) {
		this.So(this.response.StatusCode, should.Equal, http.StatusOK)
	}
	this.So(this.err, should.BeNil)
	this.So(this.fakeClient.calls, should.Equal, maxAttempts)
	this.assertBodySentToStdOut()
}

func (this *GetRetryClientFixture) assertBodySentToStdOut() {
	this.So(this.stdout.String(), should.ContainSubstring, "Internal Server Error")
}

var (
	GetRetry_ServerErrorResponse = &http.Response{StatusCode: 500, Body: newFakeBody("Internal Server Error")}
)

/////////////////////////////////////////////////////////

type FakeHTTPClientForGetRetry struct {
	calls      int
	statusCode int
}

func (this *FakeHTTPClientForGetRetry) Do(request *http.Request) (*http.Response, error) {

	this.calls++
	if request.URL.Path == "/fail-first" && this.calls < maxAttempts {
		return nil, errors.New("GOPHERS!")
	} else if request.URL.Path == "/fail-always" {
		return nil, errors.New("GOPHERS!")
	} else if request.URL.Path == "/bad-status" && this.calls < maxAttempts {
		return GetRetry_ServerErrorResponse, nil
	} else {
		return &http.Response{StatusCode: this.statusCode}, nil
	}
}

////////////////////////////////////////////////////////////////////
