package httpx

import (
	"net/http"

	"github.com/smartystreets/pipeline/numeric"
)

func WriteResponse(response http.ResponseWriter, err error) {
	if err != nil {
		writeErrorMessage(response, err.Error(), http.StatusInternalServerError)
	} else {
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}
func WriteErrorMessage(response http.ResponseWriter, message string, statusCode int) {
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(response, message, statusCode)
}

func ReadUint64Header(request *http.Request, header string) uint64 {
	raw := request.Header.Get(header)
	return numeric.StringToUint64(raw)
}
