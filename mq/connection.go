package mq

import (
	"log"

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
	log.Println("[INFO] Opening channel on existing AMQP connection.")
	if channel, err := this.inner.Channel(); err != nil {
		log.Printf("[WARN] Unable to open AMQP channel [%s]\n", err)
		return nil, err
	} else {
		log.Println("[INFO] AMQP channel opened.")
		return newChannel(channel), nil
	}
}

func (this *Connection) Close() error {
	log.Println("[INFO] Closing AMQP connection.")
	return this.inner.Close()
}
