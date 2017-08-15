package transports

import (
	"errors"
	"io"
)

type ChannelSelectWriter struct {
	input  chan []byte
	output io.Writer
}

func NewChannelSelectWriter(writer io.Writer, capacity int) *ChannelSelectWriter {
	return &ChannelSelectWriter{
		input:  make(chan []byte, capacity),
		output: writer,
	}
}

func (this *ChannelSelectWriter) Listen() {
	for buffer := range this.input {
		this.output.Write(buffer)
	}
}

func (this *ChannelSelectWriter) Write(buffer []byte) (int, error) {
	length := len(buffer)
	if length == 0 {
		return 0, nil
	}

	select {
	case this.input <- buffer:
		return length, nil
	default:
		return 0, WriteDiscardedError
	}
}

var WriteDiscardedError = errors.New("The write was discarded because the channel was full.")
