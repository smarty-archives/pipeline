package listeners

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type ShutdownListener struct {
	channel  chan os.Signal
	shutdown func()
}

func NewShutdownListener(shutdown func()) *ShutdownListener {
	channel := make(chan os.Signal, 2)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	return &ShutdownListener{channel: channel, shutdown: shutdown}
}

func (this *ShutdownListener) Listen() {
	<-this.channel
	log.Println("[INFO] Received OS shutdown signal")
	this.shutdown()
}
