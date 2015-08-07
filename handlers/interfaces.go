package pipeline

import "github.com/smartystreets/pipeline/messaging"

type Deserializer interface {
	Deserialize(*messaging.Delivery)
}
