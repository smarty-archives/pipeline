package builders

import (
	"log"
	"reflect"
	"strings"

	"github.com/smartystreets/pipeline/handlers"
	"github.com/smartystreets/pipeline/listeners"
	"github.com/smartystreets/pipeline/messaging"
)

type CompositeReaderBuilder struct {
	sourceQueue  string
	bindings     []string
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

func (this *CompositeReaderBuilder) RegisterMap(types map[string]reflect.Type) *CompositeReaderBuilder {
	for key, value := range types {
		this.types[key] = value
		this.addBinding(key)
	}

	return this
}

func (this *CompositeReaderBuilder) RegisterMultiple(prefix string, instances ...interface{}) *CompositeReaderBuilder {
	discovery := messaging.NewReflectionDiscovery(prefix)

	for _, instance := range instances {
		if discovered, err := discovery.Discover(instance); err != nil {
			log.Fatal("Unable to discover type for instance", instance)
		} else {
			this.types[discovered] = reflect.TypeOf(instance)
			this.addBinding(discovered)
		}
	}

	return this
}
func (this *CompositeReaderBuilder) Register(typeName string, instance interface{}) *CompositeReaderBuilder {
	this.types[typeName] = reflect.TypeOf(instance)
	this.addBinding(typeName)
	return this
}
func (this *CompositeReaderBuilder) addBinding(typeName string) {
	if strings.Contains(typeName, " ") {
		return // can't register .NET types
	}

	typeName = strings.Replace(typeName, ".", "-", -1)
	typeName = strings.ToLower(typeName)
	this.bindings = append(this.bindings, typeName)
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
	receive := this.openReader()
	input := receive.Deliveries()
	output := make(chan messaging.Delivery, cap(input))

	deserializer := handlers.NewJSONDeserializer(this.types)
	deserialize := handlers.NewDeserializationHandler(input, output, deserializer)
	if this.panicMissing {
		deserializer.PanicWhenMessageTypeIsUnknown()
	}
	if this.panicFail {
		deserializer.PanicWhenDeserializationFails()
	}

	return &compositeReader{
		receive:     receive,
		deserialize: deserialize,
		deliveries:  output,
	}
}

func (this *CompositeReaderBuilder) openReader() messaging.Reader {
	if len(this.sourceQueue) > 0 {
		return this.broker.OpenReader(this.sourceQueue, this.bindings...)
	}

	if len(this.bindings) > 0 {
		return this.broker.OpenTransientReader(this.bindings)
	}

	log.Fatal("Unable to open reader. No source queue or bindings specified.")
	return nil
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
