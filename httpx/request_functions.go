package httpx

import (
	"net"
	"net/http"
	"sync"

	"github.com/smartystreets/pipeline/numeric"
)

func ReadClientIPAddress(request *http.Request) string {
	if origin := request.Header.Get("X-Forwarded-For"); len(origin) > 0 {
		return origin
	} else if address, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		return address
	} else {
		return request.RemoteAddr
	}
}

func ReadUint64Header(request *http.Request, name string) uint64 {
	return numeric.StringToUint64(request.Header.Get(name))
}

func NewWaitGroup(workers int) *sync.WaitGroup {
	waiter := &sync.WaitGroup{}
	waiter.Add(workers)
	return waiter
}

func ReadHeader(request *http.Request, canoicalHeaderName string) string {
	if values, contains := request.Header[canoicalHeaderName]; contains && len(values) > 0 {
		return values[0]
	} else {
		return ""
	}
}

func WriteHeader(request *http.Request, canonicalHeaderName string, value string) {
	request.Header[canonicalHeaderName] = []string{value}
}
