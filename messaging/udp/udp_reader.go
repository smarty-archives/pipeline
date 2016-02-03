package udp

import (
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/smartystreets/clock"
	"github.com/smartystreets/logging"
	"github.com/smartystreets/pipeline/messaging"
)

type DatagramReader struct {
	acknowledgements chan interface{}
	deliveries       chan messaging.Delivery
	address          string
	mutex            *sync.Mutex
	listener         *net.UDPConn
	clock            *clock.Clock
	logger           *logging.Logger
	closed           bool
}

func NewDatagramReader(address string, capacity int) *DatagramReader {
	return &DatagramReader{
		acknowledgements: make(chan interface{}, 64),
		deliveries:       make(chan messaging.Delivery, capacity),
		address:          address,
		mutex:            &sync.Mutex{},
	}
}

func (this *DatagramReader) Listen() {
	go this.acknowledge()
	this.listen()
	close(this.acknowledgements)
}
func (this *DatagramReader) acknowledge() {
	for range this.acknowledgements {
	}
}
func (this *DatagramReader) listen() {
	this.resolveListener()

	if this.listener == nil {
		return
	}

	this.listener.SetReadBuffer(readBufferSize)

	for {
		buffer := make([]byte, readBufferSize)
		deadline := this.clock.UTCNow().Add(readDeadlineDuration)
		this.listener.SetReadDeadline(deadline)
		if read, err := this.listener.Read(buffer); err == nil {
			this.deliveries <- messaging.Delivery{
				Timestamp: this.clock.UTCNow(),
				Payload:   buffer[0:read],
			}
		} else if isClosed(err) {
			break
		}
	}
}
func (this *DatagramReader) resolveListener() {
	var listener *net.UDPConn

	for listener == nil {
		if listener, err := this.openListener(); listener != nil {
			this.logger.Printf("[INFO] Listening for UDP datagrams on %s.\n", this.address)
			this.saveListener(listener)
			break
		} else if this.isClosed() {
			break
		} else {
			this.logger.Println("[WARN] UDP socket bind failure:", err)
			time.Sleep(time.Second)
		}
	}
}
func (this *DatagramReader) openListener() (*net.UDPConn, error) {
	if address, err := net.ResolveUDPAddr("udp", this.address); err != nil {
		return nil, err
	} else if listener, err := net.ListenUDP("udp", address); err != nil {
		return nil, err
	} else {
		return listener, nil
	}
}
func (this *DatagramReader) saveListener(listener *net.UDPConn) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.listener = listener

	if this.closed {
		this.listener.Close()
	}
}
func (this *DatagramReader) isClosed() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.closed
}
func isClosed(err error) bool {
	return err == io.EOF || strings.Contains(err.Error(), "use of closed network connection")
}

func (this *DatagramReader) Close() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.closed = true

	listener := this.listener
	if listener != nil {
		listener.Close()
	}
}

func (this *DatagramReader) Deliveries() <-chan messaging.Delivery { return this.deliveries }
func (this *DatagramReader) Acknowledgements() chan<- interface{}  { return this.acknowledgements }

const (
	readDeadlineDuration = time.Second
	readBufferSize       = 64 * 1024 // 64K bytes
)
