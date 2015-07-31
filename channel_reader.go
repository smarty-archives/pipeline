package rabbit

type ChannelReader struct {
	controller       Controller
	queue            string
	bindings         []string
	acknowledgements chan interface{}
	deliveries       chan Delivery
}

func newReader(controller Controller, queue string, bindings []string) *ChannelReader {
	return &ChannelReader{
		controller:       controller,
		queue:            queue,
		bindings:         bindings,
		acknowledgements: make(chan interface{}, 1024),
		deliveries:       make(chan Delivery, 1024),
	}
}

func (this *ChannelReader) Listen() {
}

func (this *ChannelReader) Close() {
}

func (this *ChannelReader) Deliveries() <-chan Delivery {
	return this.deliveries
}
func (this *ChannelReader) Acknowledgements() chan<- interface{} {
	return this.acknowledgements
}
