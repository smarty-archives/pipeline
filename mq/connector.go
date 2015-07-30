package mq

import (
	"crypto/tls"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/smartystreets/go-rabbit"
	"github.com/streadway/amqp"
)

type Connector struct{}

func NewConnector() *Connector {
	return &Connector{}
}

func (this *Connector) Connect(target url.URL) (rabbit.Connection, error) {
	config := amqp.Config{
		TLSClientConfig: buildTLS(target.Host),
		Heartbeat:       time.Second * 15,
		Dial:            dial,
	}
	if connection, err := amqp.DialConfig(target.String(), config); err != nil {
		return nil, err
	} else {
		return newConnection(connection), nil
	}
}
func buildTLS(host string) *tls.Config {
	// FUTURE: customize TLS, e.g. acceptable list of ciphers, etc.
	return &tls.Config{
		ServerName: strings.Split(host, ":")[0],
		MinVersion: tls.VersionTLS12,
	}
}
func dial(network, address string) (net.Conn, error) {
	return net.DialTimeout(network, address, timeout)
}

var timeout = time.Second * 5 // FUTURE: customize timeout
