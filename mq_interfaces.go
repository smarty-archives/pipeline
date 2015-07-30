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
	ConfigureChannelBuffer(int) error
	ConfigureChannelAsTransactional() error

	DeclareTransientQueue() (string, error)
	BindExchangeToQueue(string, string) error

	Close() error

	Consume(string, string) (<-chan amqp.Delivery, error)
	ConsumeWithoutAcknowledgement(string, string) (<-chan amqp.Delivery, error)
	ExclusiveConsume(string, string) (<-chan amqp.Delivery, error)
	ExclusiveConsumeWithoutAcknowledgement(string, string) (<-chan amqp.Delivery, error)

	CancelConsumer(string) error

	AcknowledgeSingleMessage(uint64) error
	AcknowledgeMultipleMessages(uint64) error

	PublishMessage(string, amqp.Publishing) error

	CommitTransaction() error
	RollbackTransaction() error
}
