package persist

import (
	"net/http"
	"net/url"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type S3ClientFixture struct {
	*gunit.Fixture
	innerClient *FakeHTTPClientForSigning
	client      *S3Client
	request     *http.Request
}

const awsID = "AWSACCESSKEY"

func (this *S3ClientFixture) Setup() {
	this.innerClient = &FakeHTTPClientForSigning{}
	targetAddress, _ := url.Parse("https://s3-us-west-1.amazonaws.com/my-bucket/usage/documents/")
	this.client = NewS3Client(*targetAddress, this.innerClient, awsID, "blahblah")
	this.request, _ = http.NewRequest("PUT", "/path/to/doc.json", nil)
}

func (this *S3ClientFixture) TestRequestSigning() {
	this.client.Do(this.request)
	received := this.innerClient.received
	this.So(received, should.Equal, this.request)
	this.So(received.Host, should.Equal, "s3-us-west-1.amazonaws.com")
	this.So(received.URL.String(), should.Equal, "https://s3-us-west-1.amazonaws.com/my-bucket/usage/documents/path/to/doc.json")
	this.So(received.Header.Get("x-amz-server-side-encryption"), should.Equal, "AES256")
	this.So(received.Header.Get("Authorization"), should.StartWith, "AWS "+awsID) // ie. "AWS AWSACCESSKEY:f150Ju+KU86sySHj6pjdKQlkFhQ="
	this.So(len(received.Header.Get("Authorization")), should.BeGreaterThan, len("AWS "+awsID))
}

func (this *S3ClientFixture) TestEncryptionDisabledForGETRequests() {
	this.request.Method = "GET"
	this.client.Do(this.request)
	value, found := this.innerClient.received.Header["x-amz-server-side-encryption"]
	this.So(found, should.BeFalse)
	this.So(value, should.BeEmpty)
}

///////////////////////////////////////////////////////////////////

type FakeHTTPClientForSigning struct{ received *http.Request }

func (this *FakeHTTPClientForSigning) Do(request *http.Request) (*http.Response, error) {
	this.received = request
	return nil, nil
}
