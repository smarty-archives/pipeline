package rabbit

import "github.com/smartystreets/pipeline/messaging"

type Controller interface {
	openChannel(func() bool) Channel
	removeReader(messaging.Reader)
	removeWriter(messaging.Writer)
}
