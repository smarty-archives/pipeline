package wireup

import (
	"time"

	"github.com/smartystreets/pipeline/messaging"
)

type CompositeWriterBuilder struct {
	broker     messaging.MessageBroker
	discovery  messaging.TypeDiscovery
	retrySleep time.Duration
	panicFail  bool
}

func NewCompositeWriter(broker messaging.MessageBroker) *CompositeWriterBuilder {
	return &CompositeWriterBuilder{
		broker:     broker,
		retrySleep: time.Second * 5,
	}
}

func (this *CompositeWriterBuilder) PrefixTypesWith(prefix string) *CompositeWriterBuilder {
	this.discovery = messaging.NewReflectionDiscovery(prefix)
	return this
}

func (this *CompositeWriterBuilder) RetryAfter(sleep time.Duration) *CompositeWriterBuilder {
	this.retrySleep = sleep
	return this
}

func (this *CompositeWriterBuilder) PanicWhenSerializationFails() *CompositeWriterBuilder {
	this.panicFail = true
	return this
}

func (this *CompositeWriterBuilder) Build() messaging.CommitWriter {
	writer := this.broker.OpenTransactionalWriter()
	writer = this.layerRetry(writer)
	writer = this.layerDeserialize(writer)
	writer = this.layerDispatch(writer)
	return writer
}

func (this *CompositeWriterBuilder) layerRetry(inner messaging.CommitWriter) messaging.CommitWriter {
	if this.retrySleep == 0 {
		return inner
	}

	return messaging.NewRetryCommitWriter(inner, 0, func(uint64) {
		time.Sleep(this.retrySleep)
	})
}

func (this *CompositeWriterBuilder) layerDeserialize(inner messaging.CommitWriter) messaging.CommitWriter {
	if this.discovery == nil {
		return inner
	}

	serializer := messaging.NewJSONSerializer()
	if this.panicFail {
		serializer.PanicWhenSerializationFails()
	}

	return messaging.NewSerializationWriter(inner, serializer, this.discovery)
}

func (this *CompositeWriterBuilder) layerDispatch(inner messaging.CommitWriter) messaging.CommitWriter {
	if this.discovery == nil {
		return inner
	}

	// TODO: register overwrites and template...
	return messaging.NewDispatchWriter(inner, this.discovery)
}
