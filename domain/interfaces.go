package domain

type Handler interface {
	Handle(interface{}) []interface{}
}

type Applicator interface {
	Apply(interface{}) (modified bool)
}
