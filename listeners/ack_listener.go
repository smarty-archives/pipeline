package listeners

type AckListener struct {
	input <-chan interface{}
}

func NewAckListener(input <-chan interface{}) *AckListener {
	return &AckListener{input: input}
}

func (this *AckListener) Listen() {
	for ack := range this.input {
		if waiter, ok := ack.(WaitGroup); ok {
			waiter.Done()
		}
	}
}
