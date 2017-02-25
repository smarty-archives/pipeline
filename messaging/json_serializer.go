package messaging

import (
	"encoding/json"

	"github.com/smartystreets/logging"
)

type JSONSerializer struct {
	logger *logging.Logger

	panicFail bool
}

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

func (this *JSONSerializer) PanicWhenSerializationFails() {
	this.panicFail = true
}

func (this *JSONSerializer) Serialize(message interface{}) ([]byte, error) {
	content, err := json.Marshal(message)
	if this.panicFail && err != nil {
		this.logger.Panic("[ERR] Could not serialize message:", err)
	} else if err != nil {
		this.logger.Println("[ERR] Could not serialize message:", err)
	}
	return content, err
}

func (this *JSONSerializer) ContentType() string     { return "application/json" }
func (this *JSONSerializer) ContentEncoding() string { return "" }
