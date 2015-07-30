package mq

import (
	"github.com/smartystreets/go-rabbit"
	"github.com/streadway/amqp"
)

type RabbitConnection struct {
	inner *amqp.Connection
}

func newConnection(inner *amqp.Connection) rabbit.Connection {
	return &RabbitConnection{inner: inner}
}

func (this *RabbitConnection) Channel() (rabbit.Channel, error) {
	if channel, err := this.inner.Channel(); err != nil {
		return nil, err
	} else {
		return newChannel(channel), nil
	}
}

func (this *RabbitConnection) Close() error {
	return this.inner.Close()
}
