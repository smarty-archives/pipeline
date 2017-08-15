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
	this.events = this.events[0:0]

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

func (this *SimpleMessageRouter) Add(item interface{}) bool {
	added := false

	if handler, ok := item.(Handler); ok {
		this.AddHandler(handler)
		added = true
	}

	if applicator, ok := item.(Applicator); ok {
		this.AddDocument(applicator)
		added = true
	}

	return added
}

func (this *SimpleMessageRouter) AddHandler(handler Handler) {
	this.handlers = append(this.handlers, handler)
}

func (this *SimpleMessageRouter) AddDocument(document Applicator) {
	this.documents = append(this.documents, document)
}
