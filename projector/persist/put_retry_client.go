package persist

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type PutRetryClient struct {
	inner   HTTPClient
	retries int
}

func NewPutRetryClient(inner HTTPClient, retries int) *PutRetryClient {
	return &PutRetryClient{inner: inner, retries: retries}
}

func (this *PutRetryClient) Do(request *http.Request) (*http.Response, error) {
	request.Body = newRetryBuffer(request.Body)

	for current := 0; current <= this.retries; current++ {
		log.Printf("[INFO] Saving document: '%s' (attempt #%d)\n", request.URL.Path, current)

		response, err := this.inner.Do(request)

		if err == nil && response.StatusCode == http.StatusOK {
			log.Println("[INFO] Document saved:", request.URL.Path)
			return response, nil
		} else if err != nil {
			log.Println("[WARN] Unexpected response from target storage:", err)
		} else if response.Body != nil {
			log.Printf("[WARN] Target host rejected request ('%s'):\n%s\n", request.URL.Path, readResponse(response))
		}

		napTime(time.Second * 10)
	}

	return nil, errors.New("Max retries exceeded. Unable to connect.")
}
func readResponse(response *http.Response) string {
	responseDump, _ := httputil.DumpResponse(response, true)
	return string(responseDump) + "\n-------------------------------------------"
}

var napTime = func(duration time.Duration) {
	time.Sleep(duration)
}

type retryBuffer struct{ io.ReadSeeker }

func newRetryBuffer(body io.ReadCloser) *retryBuffer {
	return &retryBuffer{body.(io.ReadSeeker)}
}
func (this *retryBuffer) Close() error {
	this.Seek(0, 0) // seeks to the beginning (to allow retry) when the buffer is "Closed"
	return nil
}
