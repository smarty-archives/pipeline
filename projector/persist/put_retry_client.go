package persist

import (
	"errors"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/smartystreets/clock"
	"github.com/smartystreets/logging"
)

type PutRetryClient struct {
	inner   HTTPClient
	retries int
	sleeper *clock.Sleeper
	logger  *logging.Logger
}

func NewPutRetryClient(inner HTTPClient, retries int) *PutRetryClient {
	return &PutRetryClient{inner: inner, retries: retries}
}

func (this *PutRetryClient) Do(request *http.Request) (*http.Response, error) {
	request.Body = newRetryBuffer(request.Body)

	for current := 0; current <= this.retries; current++ {
		response, err := this.inner.Do(request)

		if err == nil && response.StatusCode == http.StatusOK {
			return response, nil
		} else if err != nil {
			this.logger.Println("[WARN] Unexpected response from target storage:", err)
		} else if response.Body != nil {
			this.logger.Printf("[WARN] Target host rejected request ('%s'):\n%s\n", request.URL.Path, readResponse(response))
		}

		this.sleeper.Sleep(time.Second * 10)
	}

	return nil, errors.New("Max retries exceeded. Unable to connect.")
}
func readResponse(response *http.Response) string {
	responseDump, _ := httputil.DumpResponse(response, true)
	return string(responseDump) + "\n-------------------------------------------"
}

type retryBuffer struct{ io.ReadSeeker }

func newRetryBuffer(body io.ReadCloser) *retryBuffer {
	return &retryBuffer{body.(io.ReadSeeker)}
}
func (this *retryBuffer) Close() error {
	this.Seek(0, 0) // seeks to the beginning (to allow retry) when the buffer is "Closed"
	return nil
}
