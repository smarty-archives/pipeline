package pipeline

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/smartystreets/go-messenger"
)

type JSONDeserializer struct {
	types map[string]reflect.Type

	panicMissingType bool
	panicUnmarshal   bool
}

func NewJSONDeserializer(types map[string]reflect.Type) *JSONDeserializer {
	return &JSONDeserializer{types: types}
}

func (this *JSONDeserializer) PanicWhenMessageTypeIsUnkonwn() {
	this.panicMissingType = true
}

func (this *JSONDeserializer) PanicWhenDeserializationFails() {
	this.panicUnmarshal = true
}

func (this *JSONDeserializer) Deserialize(delivery *messenger.Delivery) {
	messageType, found := this.types[delivery.MessageType]
	if !found && this.panicMissingType {
		log.Panicf("MessageType not found: '%s'", delivery.MessageType)
	} else if !found {
		return
	}

	message := reflect.New(messageType).Interface()
	err := json.Unmarshal(delivery.Payload, message)
	if err != nil && this.panicUnmarshal {
		log.Panicf("Could not deserialize message of type '%s': %s", delivery.MessageType, err.Error())
	} else if err != nil {
		return
	}
	delivery.Message = message
}
