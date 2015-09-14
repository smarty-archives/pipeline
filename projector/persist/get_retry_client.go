package persist

import (
	"errors"
	"log"
	"net/http"
	"time"
)

type GetRetryClient struct {
	inner   HTTPClient
	retries int
	napTime func(time.Duration)
}

func NewGetRetryClient(inner HTTPClient, retries int) *GetRetryClient {
	return &GetRetryClient{inner: inner, retries: retries, napTime: napTime}
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
		this.napTime(time.Second * 5)
	}
	return nil, errors.New("Max retries exceeded. Unable to connect.")
}
