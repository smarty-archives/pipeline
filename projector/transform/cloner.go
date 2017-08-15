package transform

import (
	"encoding/gob"
	"reflect"

	"github.com/smartystreets/logging"
	"github.com/smartystreets/pipeline/projector"
)

type DocumentCloner struct {
	logger *logging.Logger

	buffer  ResetReadWriter
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func NewDocumentCloner(buffer ResetReadWriter) *DocumentCloner {
	return &DocumentCloner{
		buffer:  buffer,
		encoder: gob.NewEncoder(buffer),
		decoder: gob.NewDecoder(buffer),
	}
}

func (this *DocumentCloner) Clone(document projector.Document) projector.Document {
	this.buffer.Reset()
	this.encode(document)
	clone := blank(document)
	this.decode(clone)
	return clone
}

func (this *DocumentCloner) encode(document projector.Document) {
	if err := this.encoder.Encode(document); err != nil {
		this.logger.Panic("Gob encoding error:", err)
	}
}
func (this *DocumentCloner) decode(clone projector.Document) {
	if err := this.decoder.Decode(clone); err != nil {
		this.logger.Panic("Gob decoding error:", err)
	}
}

// http://stackoverflow.com/questions/7850140/how-do-you-create-a-new-instance-of-a-struct-from-its-type-at-runtime-in-go
func blank(document projector.Document) projector.Document {
	return reflect.New(reflect.TypeOf(document).Elem()).Interface().(projector.Document)
}
