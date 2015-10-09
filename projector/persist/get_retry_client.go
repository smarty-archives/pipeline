package persist

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/smartystreets/clock"
)

type GetRetryClient struct {
	inner   HTTPClient
	retries int
	sleeper *clock.Sleeper
}

// FUTURE: We may want to consider a ShutdownClient that sits just under
// the RetryClient. This makes it possible for a shutdown signal to break
// a retry loop because the Shutdown client would retry success (HTTP 200)
// or perhaps HTTP 404?

func NewGetRetryClient(inner HTTPClient, retries int) *GetRetryClient {
	return &GetRetryClient{inner: inner, retries: retries}
}

func (this *GetRetryClient) Do(request *http.Request) (*http.Response, error) {
	for current := 0; current <= this.retries; current++ {
		response, err := this.inner.Do(request)
		if err == nil && (response.StatusCode == http.StatusOK || response.StatusCode == http.StatusNotFound) {
			return response, nil
		} else if err != nil {
			log.Println("[WARN] Unexpected response from target storage:", err)
		} else if response.Body != nil {
			log.Printf("[WARN] Target host rejected request ('%s'):\n%s\n", request.URL.Path, readResponse(response))
		}
		this.sleeper.Sleep(time.Second * 5)
	}
	return nil, errors.New("Max retries exceeded. Unable to connect.")
}
