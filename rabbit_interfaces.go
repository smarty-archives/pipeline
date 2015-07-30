package rabbit

type (
	Broker interface {
		Connect() error
		Disconnect()

		OpenReader(queue string) Reader
		OpenTransientReader(bindings []string) Reader

		OpenWriter() Writer
		OpenTransactionalWriter() CommitWriter
	}

	Reader interface {
		Listen()
		Close()

		Deliveries() <-chan Delivery
		Acknowledgements() chan<- interface{}
	}

	Writer interface {
		Write(Dispatch) error
		Close()
	}

	CommitWriter interface {
		Writer
	}
)

type Controller interface {
	openChannel() Channel
	removeReader(Reader)
	removeWriter(Writer)
}
