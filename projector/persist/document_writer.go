package persist

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/smartystreets/logging"
	"github.com/smartystreets/pipeline/projector"
)

type DocumentWriter struct {
	logger *logging.Logger

	client HTTPClient
}

func NewDocumentWriter(client HTTPClient) *DocumentWriter {
	return &DocumentWriter{client: client}
}

func (this *DocumentWriter) Write(document projector.Document) {
	body := this.serialize(document)
	checksum := this.md5Checksum(body)
	request := this.buildRequest(document.Path(), body, checksum)
	response, err := this.client.Do(request)
	this.handleResponse(response, err)
}

func (this *DocumentWriter) serialize(document projector.Document) []byte {
	buffer := bytes.NewBuffer([]byte{})
	gzipper, _ := gzip.NewWriterLevel(buffer, gzip.BestCompression)
	encoder := json.NewEncoder(gzipper)

	if err := encoder.Encode(document); err != nil {
		this.logger.Panic(err)
	}

	gzipper.Close()
	return buffer.Bytes()
}

func (this *DocumentWriter) md5Checksum(body []byte) string {
	sum := md5.Sum(body)
	return base64.StdEncoding.EncodeToString(sum[:])
}

func (this *DocumentWriter) buildRequest(path string, body []byte, checksum string) *http.Request {
	request, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		this.logger.Panic(err)
	}

	request.Body = newNopCloser(body)
	request.ContentLength = int64(len(body))

	this.setHeaders(request, checksum)
	return request
}
func (this *DocumentWriter) setHeaders(request *http.Request, checksum string) {
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Content-Md5", checksum)
	request.Header.Set("Content-Type", "application/json")
}

// handleResponse handles error response, which technically, shouldn't happen
// because the inner client should be handling retry indefinitely, until the service
// response. This is here merely for the sake of completeness, and to bullet-proof
// the software in case the behavior of the inner client changes in the future.
func (this *DocumentWriter) handleResponse(response *http.Response, err error) {
	if err != nil {
		this.logger.Panic(err)
	} else if response.StatusCode != http.StatusOK {
		this.logger.Panic(fmt.Errorf("Non-200 HTTP Status Code: %d %s", response.StatusCode, response.Status))
	}

	if response != nil && response.Body != nil {
		response.Body.Close() // release connection back to pool
	}
}

type nopCloser struct{ io.ReadSeeker }

func newNopCloser(body []byte) *nopCloser { return &nopCloser{bytes.NewReader(body)} }
func (this *nopCloser) Close() error      { return nil }
