package rabbit

import "github.com/smartystreets/go-messenger"

type Controller interface {
	openChannel() Channel
	removeReader(messenger.Reader)
	removeWriter(messenger.Writer)
}
