package persist

import (
	"net/http"
	"net/url"
	"path"

	"github.com/smartystreets/go-aws-auth"
)

type S3Client struct {
	target      url.URL
	inner       HTTPClient
	credentials awsauth.Credentials
}

func NewS3Client(target url.URL, inner HTTPClient, awsAccessID, awsSecretKey string) *S3Client {
	return &S3Client{
		target: target,
		inner:  inner,
		credentials: awsauth.Credentials{
			AccessKeyID:     awsAccessID,
			SecretAccessKey: awsSecretKey,
		},
	}
}

func (this *S3Client) Do(request *http.Request) (*http.Response, error) {
	raw := "https://" + this.target.Host + path.Join(this.target.Path, request.URL.Path)
	request.URL, _ = url.Parse(raw) // URL previously validated
	request.Host = this.target.Host
	if request.Method != "GET" {
		request.Header.Set("x-amz-server-side-encryption", "AES256")
	}

	// don't use Sign4, it reads and replaces the body which affects retry
	awsauth.SignS3(request, this.credentials)

	return this.inner.Do(request)
}
