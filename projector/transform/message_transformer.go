package transform

import (
	"time"

	"github.com/smartystreets/pipeline/projector"
)

type MessageTransformer struct {
	documents []projector.Document
	changed   map[string]projector.Document
	cloner    Cloner
}

func NewMessageTransformer(documents []projector.Document, cloner Cloner) *MessageTransformer {
	return &MessageTransformer{
		documents: documents,
		changed:   make(map[string]projector.Document, 16),
		cloner:    cloner,
	}
}

func (this *MessageTransformer) TransformAllDocuments(message interface{}, now time.Time) {
	if message == nil {
		return
	}

	for i, doc := range this.documents {
		doc = doc.Lapse(now)
		this.documents[i] = doc
		if doc.Apply(message) {
			this.changed[doc.Path()] = doc
		}
	}
}

func (this *MessageTransformer) Collect() []projector.Document {
	docs := []projector.Document{}

	for key, doc := range this.changed {
		delete(this.changed, key)
		docs = append(docs, this.cloner.Clone(doc))
	}

	return docs
}
