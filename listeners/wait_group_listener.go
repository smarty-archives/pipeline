package listeners

type WaitGroupListener struct {
	inner  Listener
	waiter WaitGroup
}

func NewWaitGroupListener(listener Listener, waiter WaitGroup) *WaitGroupListener {
	return &WaitGroupListener{
		inner:  listener,
		waiter: waiter,
	}
}

func (this *WaitGroupListener) Listen() {
	if this.inner == nil {
		return
	}

	this.waiter.Add(1)
	this.inner.Listen()
	this.waiter.Done()
}

func (this *WaitGroupListener) Close() {
	if closer, ok := this.inner.(ListenCloser); ok {
		closer.Close()
	}
}
