package listeners

import "sync"

type CompositeWaitListener struct {
	waiter *sync.WaitGroup
	inner  Listener
}

func NewCompositeWaitShutdownListener(primary ListenCloser, listeners ...Listener) *CompositeWaitListener {
	listeners = append(listeners, primary, NewShutdownListener(primary.Close))
	return NewCompositeWaitListener(listeners...)
}
func NewCompositeWaitListener(listeners ...Listener) *CompositeWaitListener {
	waiter := &sync.WaitGroup{}

	items := make([]Listener, 0, len(listeners))
	for _, item := range listeners {
		item = NewWaitGroupListener(item, waiter)
		items = append(items, item)
	}

	return &CompositeWaitListener{
		waiter: waiter,
		inner:  NewCompositeListener(items),
	}
}

func (this *CompositeWaitListener) Listen() {
	this.inner.Listen()
	this.waiter.Wait()
}
