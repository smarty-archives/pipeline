package rabbit

import (
	"net/url"

	"github.com/streadway/amqp"
)

type Connector interface {
	Connect(url.URL) (Connection, error)
}

type Connection interface {
	Channel() (Channel, error)
	Close() error
}

type Channel interface {
	Consumer
	Publisher

	Close() error
}

type Consumer interface {
	Acknowledger

	ConfigureChannelBuffer(int) error
	DeclareQueue(string) error
	DeclareTransientQueue() (string, error)
	BindExchangeToQueue(string, string) error

	Consume(string, string) (<-chan amqp.Delivery, error)
	ExclusiveConsume(string, string) (<-chan amqp.Delivery, error)
	ConsumeWithoutAcknowledgement(string, string) (<-chan amqp.Delivery, error)
	ExclusiveConsumeWithoutAcknowledgement(string, string) (<-chan amqp.Delivery, error)

	CancelConsumer(string) error
}

type Acknowledger interface {
	AcknowledgeSingleMessage(uint64) error
	AcknowledgeMultipleMessages(uint64) error
}

type Publisher interface {
	ConfigureChannelAsTransactional() error

	PublishMessage(string, amqp.Publishing) error

	CommitTransaction() error
	RollbackTransaction() error
}
