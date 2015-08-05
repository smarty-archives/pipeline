package messenger

import "errors"

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
		Commit() error
	}
)

var WriterClosedError error = errors.New("The writer has been closed and can no longer be used.")
