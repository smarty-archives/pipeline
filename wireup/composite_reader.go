package wireup

import (
	"log"
	"reflect"

	"github.com/smartystreets/pipeline/handlers"
	"github.com/smartystreets/pipeline/listeners"
	"github.com/smartystreets/pipeline/messaging"
)

type CompositeReaderBuilder struct {
	sourceQueue  string
	broker       messaging.MessageBroker
	types        map[string]reflect.Type
	panicMissing bool
	panicFail    bool
}

func NewCompositeReader(broker messaging.MessageBroker, sourceQueue string) *CompositeReaderBuilder {
	return &CompositeReaderBuilder{
		broker:      broker,
		sourceQueue: sourceQueue,
		types:       make(map[string]reflect.Type),
	}
}

func (this *CompositeReaderBuilder) Register(prefix string, instances ...interface{}) *CompositeReaderBuilder {
	discovery := messaging.NewReflectionDiscovery(prefix)

	for _, instance := range instances {
		if discovered, err := discovery.Discover(instance); err != nil {
			log.Fatal("Unable to discover type for instance", instance)
		} else {
			this.types[discovered] = reflect.TypeOf(instance)
		}
	}

	return this
}

func (this *CompositeReaderBuilder) RegisterType(name string, value reflect.Type) *CompositeReaderBuilder {
	this.types[name] = value
	return this
}

func (this *CompositeReaderBuilder) PanicWhenMessageTypeNotFound() *CompositeReaderBuilder {
	this.panicMissing = true
	return this
}

func (this *CompositeReaderBuilder) PanicWhenDeserializationFails() *CompositeReaderBuilder {
	this.panicFail = true
	return this
}

func (this *CompositeReaderBuilder) Build() messaging.Reader {
	receive := this.broker.OpenReader(this.sourceQueue)
	input := receive.Deliveries()
	output := make(chan messaging.Delivery, cap(input))
	deserializer := handlers.NewJSONDeserializer(this.types)
	deserialize := handlers.NewDeserializationHandler(input, output, deserializer)

	return &compositeReader{
		receive:     receive,
		deserialize: deserialize,
		deliveries:  output,
	}
}

type compositeReader struct {
	receive     messaging.Reader
	deserialize listeners.Listener
	deliveries  chan messaging.Delivery
}

func (this *compositeReader) Listen() {
	listeners.NewCompositeWaitListener(
		this.receive,
		this.deserialize,
	).Listen()
}
func (this *compositeReader) Close() {
	this.receive.Close()
}

func (this *compositeReader) Deliveries() <-chan messaging.Delivery {
	return this.deliveries
}
func (this *compositeReader) Acknowledgements() chan<- interface{} {
	return this.receive.Acknowledgements()
}
