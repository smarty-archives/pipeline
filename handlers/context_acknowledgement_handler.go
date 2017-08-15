package handlers

type ContextAcknowledgementHandler struct {
	input <-chan RequestContext
}

func NewContextAcknowledgementHandler(input <-chan RequestContext) *ContextAcknowledgementHandler {
	return &ContextAcknowledgementHandler{input: input}
}

func (this *ContextAcknowledgementHandler) Listen() {
	for context := range this.input {
		context.Close()
	}
}
