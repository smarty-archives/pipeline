package rabbit

import "github.com/smartystreets/pipeline/messaging"

type Controller interface {
	openChannel() Channel
	removeReader(messaging.Reader)
	removeWriter(messaging.Closer)
}
