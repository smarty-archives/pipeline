package transports

import (
	"log"
	"net"
	"net/url"
	"time"
)

type UDPTransport struct {
	target        *url.URL
	canLogFailure bool
	sockets       []net.Conn
	resolveDNS    bool
	index         int
}

func NewUDPTransport(target *url.URL) *UDPTransport {
	this := &UDPTransport{target: target, canLogFailure: true, index: -1}
	go this.beginResolve()
	return this
}

func (this *UDPTransport) Write(payload []byte) (int, error) {
	if err := this.ensureConnection(); err != nil {
		this.logFailure(err)
		this.closeConnection()
		return 0, err
	} else if written, err := this.write(payload); err != nil {
		this.logFailure(err)
		this.closeConnection()
		return 0, err
	} else {
		return written, nil
	}
}
func (this *UDPTransport) write(payload []byte) (int, error) {
	if length := len(this.sockets); length == 0 {
		return 0, nil
	} else if length == 1 {
		return this.sockets[0].Write(payload)
	} else {
		this.index = (this.index + 1) % length
		return this.sockets[this.index].Write(payload)
	}
}

func (this *UDPTransport) ensureConnection() error {
	if this.resolveDNS {
		this.resolveDNS = false
		this.closeConnection()
	}

	if len(this.sockets) > 0 {
		return nil
	} else if sockets, err := this.dial(); err != nil {
		return err
	} else {
		log.Println("[INFO] UDP Transport Connection(s) Established:", len(sockets))
		this.sockets = sockets
		this.canLogFailure = true
	}

	return nil
}
func (this *UDPTransport) dial() ([]net.Conn, error) {
	host, port, err := net.SplitHostPort(this.target.Host)
	if err != nil {
		return nil, err
	}

	targets, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	sockets := []net.Conn{}
	for _, target := range targets {
		address := target.String() + ":" + port
		if socket, err := net.DialTimeout("udp", address, time.Second*5); err != nil {
			return nil, err
		} else {
			sockets = append(sockets, socket)
		}
	}

	return sockets, nil
}
func (this *UDPTransport) logFailure(err error) {
	if this.canLogFailure {
		this.canLogFailure = false // don't log again until we have a good socket
		log.Println("[WARNING] UDP Transport Connection Failure:", err)
	}
}
func (this *UDPTransport) closeConnection() {
	for _, socket := range this.sockets {
		socket.Close()
	}

	this.sockets = this.sockets[0:0]
}
func (this *UDPTransport) beginResolve() {
	for {
		time.Sleep(time.Minute * 1)
		this.resolveDNS = true
	}
}
