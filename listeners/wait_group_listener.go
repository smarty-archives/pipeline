package listeners

import "sync"

type WaitGroupListener struct {
	inner  Listener
	waiter *sync.WaitGroup
}

func NewWaitGroupListener(listener Listener, waiter *sync.WaitGroup) Listener {
	waiter.Add(1)

	return &WaitGroupListener{
		inner:  listener,
		waiter: waiter,
	}
}

func (this *WaitGroupListener) Listen() {
	this.inner.Listen()
	this.waiter.Done()
}
