package builders

import (
	"time"

	"github.com/smartystreets/messaging"
)

type CompositeWriterBuilder struct {
	broker        messaging.MessageBroker
	discovery     messaging.TypeDiscovery
	retrySleep    time.Duration
	retryCallback func(uint64)
	maxRetries    uint64
	panicFail     bool
}

func NewCompositeWriter(broker messaging.MessageBroker) *CompositeWriterBuilder {
	return &CompositeWriterBuilder{
		broker:     broker,
		retrySleep: time.Second * 5,
	}
}

func (this *CompositeWriterBuilder) WithDiscovery(discovery messaging.TypeDiscovery) *CompositeWriterBuilder {
	this.discovery = discovery
	return this
}

func (this *CompositeWriterBuilder) WithMaxRetries(maxRetries uint64) *CompositeWriterBuilder {
	this.maxRetries = maxRetries
	return this
}

func (this *CompositeWriterBuilder) RetryAfterSleep(sleep time.Duration) *CompositeWriterBuilder {
	this.retrySleep = sleep
	return this
}

func (this *CompositeWriterBuilder) RetryAfterCallback(callback func(uint64)) *CompositeWriterBuilder {
	this.retryCallback = callback
	return this
}

func (this *CompositeWriterBuilder) PanicWhenSerializationFails() *CompositeWriterBuilder {
	this.panicFail = true
	return this
}

func (this *CompositeWriterBuilder) Build() messaging.CommitWriter {
	writer := this.broker.OpenTransactionalWriter()
	writer = this.layerRetry(writer)
	writer = this.layerSerialize(writer)
	writer = this.layerDispatch(writer)
	return writer
}

func (this *CompositeWriterBuilder) layerRetry(inner messaging.CommitWriter) messaging.CommitWriter {
	if this.retryCallback != nil {
		return messaging.NewRetryCommitWriter(inner, this.maxRetries, nil, this.retryCallback)
	}

	if this.retrySleep <= 0 {
		return inner
	}

	return messaging.NewRetryCommitWriter(inner, this.maxRetries, nil, func(uint64) {
		time.Sleep(this.retrySleep)
	})
}

func (this *CompositeWriterBuilder) layerSerialize(inner messaging.CommitWriter) messaging.CommitWriter {
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

	return messaging.NewDispatchWriter(inner, this.discovery)
}
