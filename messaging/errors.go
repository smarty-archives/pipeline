package messaging

import "errors"

var (
	WriterClosedError         = errors.New("The writer has been closed and can no longer be used.")
	BrokerShuttingDownError   = errors.New("Broker is still shutting down.")
	EmptyDispatchError        = errors.New("Unable to write an empty dispatch")
	MessageTypeDiscoveryError = errors.New("Unable to discover message type")
)
