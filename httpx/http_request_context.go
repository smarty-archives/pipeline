package httpx

type HTTPRequestContext struct {
	Result []interface{}
	Waiter WaitGroup
}

func NewRequestContext(waiter WaitGroup) *HTTPRequestContext {
	waiter.Add(1)

	return &HTTPRequestContext{Waiter: waiter}
}

func (this *HTTPRequestContext) Write(message interface{}) {
	this.Result = append(this.Result, message)
}

func (this *HTTPRequestContext) Written() interface{} {
	if length := len(this.Result); length == 0 {
		return nil
	} else if length == 1 {
		return this.Result[0]
	} else {
		return this.Result
	}
}

func (this *HTTPRequestContext) Close() {
	this.Waiter.Done()
}
