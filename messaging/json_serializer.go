package messaging

import (
	"encoding/json"
	"log"
)

type JSONSerializer struct {
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
		log.Panic("Could not deserialize message:", err)
	} else if err != nil {
		log.Println("Could not deserialize message:", err)
	}
	return content, err
}
