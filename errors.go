package messenger

import "errors"

var (
	WriterClosedError       = errors.New("The writer has been closed and can no longer be used.")
	BrokerShuttingDownError = errors.New("Broker is still shutting down.")
)
