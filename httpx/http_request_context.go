package httpx

type HTTPRequestContext struct {
	Result interface{}
	Waiter WaitGroup
}

func NewRequestContext(waiter WaitGroup) *HTTPRequestContext {
	waiter.Add(1)

	return &HTTPRequestContext{Waiter: waiter}
}

func (this *HTTPRequestContext) Write(message interface{}) {
	this.Result = message
}

func (this *HTTPRequestContext) Close() {
	this.Waiter.Done()
}
