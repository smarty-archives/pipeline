package domain

type SimpleMessageRouter struct {
	events    []interface{}
	handlers  []Handler
	documents []Applicator
}

func NewSimpleMessageRouter(handlers []Handler, documents []Applicator) *SimpleMessageRouter {
	return &SimpleMessageRouter{handlers: handlers, documents: documents}
}

func (this *SimpleMessageRouter) Handle(command interface{}) []interface{} {
	this.events = this.events[:]

	for _, handler := range this.handlers {
		for _, event := range handler.Handle(command) {
			this.Apply(event)
			this.events = append(this.events, event)
		}
	}

	return this.events
}

func (this *SimpleMessageRouter) Apply(event interface{}) bool {
	result := false
	for _, doc := range this.documents {
		result = doc.Apply(event) || result
	}
	return result
}

func (this *SimpleMessageRouter) AddHandler(handler Handler) {
	this.handlers = append(this.handlers, handler)
}

func (this *SimpleMessageRouter) AddDocument(document Applicator) {
	this.documents = append(this.documents, document)
}
