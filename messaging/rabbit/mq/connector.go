package mq

import (
	"crypto/tls"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/smartystreets/pipeline/messaging/rabbit"
	"github.com/streadway/amqp"
)

type Connector struct{}

func NewConnector() *Connector {
	return &Connector{}
}

func (this *Connector) Connect(target url.URL) (rabbit.Connection, error) {
	config := amqp.Config{
		Heartbeat: time.Second * 15,
		Dial:      customDialer,
	}

	log.Println("[INFO] Establishing connection to AMQP broker.")
	if connection, err := amqp.DialConfig(target.String(), config); err != nil {
		log.Println("[WARN] Unable to establish connection", err)
		return nil, err
	} else {
		log.Println("[INFO] AMQP connection established.")
		return newConnection(connection), nil
	}
}
func customDialer(network, target string) (net.Conn, error) {
	connection, err := defaultDialer(network, target)
	if err != nil {
		return nil, err
	}

	tlsConfig := buildTLS(target)
	if tlsConfig == nil {
		return connection, nil
	}

	tlsClient := tls.Client(connection, tlsConfig)
	tlsClient.SetDeadline(time.Now().Add(timeout))
	if err := tlsClient.Handshake(); err != nil {
		connection.Close()
		return nil, err
	}

	return tlsClient, nil
}

func defaultDialer(network, target string) (net.Conn, error) {
	connection, err := net.DialTimeout(network, target, timeout)
	if err != nil {
		return nil, err
	}

	// Heartbeating hasn't started yet, don't stall forever on a dead server.
	if err := connection.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	return connection, nil
}

func buildTLS(address string) *tls.Config {
	target, _ := url.Parse(address)
	if strings.ToLower(target.Scheme) != "amqps" {
		return nil
	}

	// FUTURE: customize TLS, e.g. acceptable list of ciphers, etc.
	return &tls.Config{
		ServerName: strings.Split(target.Host, ":")[0],
		MinVersion: tls.VersionTLS12,
	}
}

var timeout = time.Second * 5 // FUTURE: customize timeout
