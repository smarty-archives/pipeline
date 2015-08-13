package messaging

type (
	MessageBroker interface {
		Connect() error
		Disconnect()

		OpenReader(queue string) Reader
		OpenTransientReader(bindings []string) Reader

		OpenWriter() Writer
		OpenTransactionalWriter() CommitWriter
	}

	Reader interface {
		Listen()
		Closer

		Deliveries() <-chan Delivery
		Acknowledgements() chan<- interface{}
	}

	Writer interface {
		Write(Dispatch) error
		Closer
	}

	CommitWriter interface {
		Write(Dispatch)
		Commit() error
		Closer
	}

	Closer interface {
		Close()
	}
)
