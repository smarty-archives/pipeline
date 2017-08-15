package handlers

type EventMessage struct {
	Message    interface{}
	Context    RequestContext
	EndOfBatch bool
}
