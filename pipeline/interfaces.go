package pipeline

import "github.com/smartystreets/go-messenger"

type Deserializer interface {
	Deserialize(*messenger.Delivery)
}
