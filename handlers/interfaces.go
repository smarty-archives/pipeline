package handlers

import "github.com/smartystreets/pipeline/messaging"

type Deserializer interface {
	Deserialize(*messaging.Delivery)
}

type RequestContext interface {
	Write(interface{})
	Close()
}

type WaitGroup interface {
	Add(delta int)
	Done()
}
