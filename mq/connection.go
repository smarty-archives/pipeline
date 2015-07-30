package mq

import (
	"github.com/smartystreets/go-rabbit"
	"github.com/streadway/amqp"
)

type Connection struct {
	inner *amqp.Connection
}

func newConnection(inner *amqp.Connection) rabbit.Connection {
	return &Connection{inner: inner}
}

func (this *Connection) Channel() (rabbit.Channel, error) {
	if channel, err := this.inner.Channel(); err != nil {
		return nil, err
	} else {
		return newChannel(channel), nil
	}
}

func (this *Connection) Close() error {
	return this.inner.Close()
}
