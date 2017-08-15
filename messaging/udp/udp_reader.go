package udp

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/smartystreets/clock"
	"github.com/smartystreets/pipeline/messaging"
)

type DatagramReader struct {
	acknowledgements chan interface{}
	deliveries       chan messaging.Delivery
	socket           net.Conn
	clock            *clock.Clock
}

func NewDatagramReader(bindAddress string, capacity int) *DatagramReader {
	socket, _ := BindSocket(bindAddress)
	return NewDatagramReaderWithSocket(socket, capacity)
}
func NewDatagramReaderWithSocket(socket net.Conn, capacity int) *DatagramReader {
	return &DatagramReader{
		acknowledgements: make(chan interface{}, capacity),
		deliveries:       make(chan messaging.Delivery, capacity),
		socket:           socket,
	}
}

func (this *DatagramReader) Listen() {
	go this.acknowledge()
	this.listen()
	close(this.deliveries)
}
func (this *DatagramReader) acknowledge() {
	for range this.acknowledgements {
	}
}
func (this *DatagramReader) listen() {
	if this.socket == nil {
		return
	}

	for {
		deadline := this.clock.UTCNow().Add(readDeadlineDuration)
		this.socket.SetReadDeadline(deadline)

		buffer := make([]byte, readBufferSize)
		if read, err := this.socket.Read(buffer); err == nil {
			this.deliveries <- messaging.Delivery{
				Timestamp: this.clock.UTCNow(),
				Payload:   buffer[0:read],
			}
		} else if err == io.EOF {
			break
		} else if strings.Contains(err.Error(), "use of closed network connection") {
			break
		} else {
			continue // e.g. timeout messages
		}
	}
}

func (this *DatagramReader) Close() {
	if this.socket != nil {
		this.socket.Close()
	}
}

func (this *DatagramReader) Deliveries() <-chan messaging.Delivery { return this.deliveries }
func (this *DatagramReader) Acknowledgements() chan<- interface{}  { return this.acknowledgements }

const (
	readDeadlineDuration = time.Second
	readBufferSize       = 64 * 1024 // 64K bytes
)
