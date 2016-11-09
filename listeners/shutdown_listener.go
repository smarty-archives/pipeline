package listeners

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/smartystreets/logging"
)

type ShutdownListener struct {
	logger *logging.Logger

	mutex    sync.Once
	channel  chan os.Signal
	shutdown func()
}

func NewShutdownListener(shutdown func()) *ShutdownListener {
	channel := make(chan os.Signal, 2)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	return &ShutdownListener{channel: channel, shutdown: shutdown}
}

func (this *ShutdownListener) Listen() {
	if <-this.channel != nil {
		this.logger.Println("[INFO] Received OS shutdown signal.")
	}

	this.shutdown()
}

func (this *ShutdownListener) Close() {
	this.mutex.Do(this.close)
}

func (this *ShutdownListener) close() {
	signal.Stop(this.channel)
	close(this.channel)
	this.logger.Println("[INFO] Unsubscribed from OS shutdown signals.")
}
