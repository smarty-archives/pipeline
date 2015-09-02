package listeners

import "sync"

type CompositeWaitListener struct {
	waiter *sync.WaitGroup
	items  []Listener
}

func NewCompositeWaitShutdownListener(primary ListenCloser, listeners ...Listener) *CompositeWaitListener {
	listeners = append(listeners, primary, NewShutdownListener(primary.Close))
	return NewCompositeWaitListener(listeners...)
}

func NewCompositeWaitListener(listeners ...Listener) *CompositeWaitListener {
	return &CompositeWaitListener{
		waiter: &sync.WaitGroup{},
		items:  listeners,
	}
}

func (this *CompositeWaitListener) Listen() {
	this.waiter.Add(len(this.items))

	for _, item := range this.items {
		go this.listen(item)
	}

	this.waiter.Wait()
}

func (this *CompositeWaitListener) listen(listener Listener) {
	if listener != nil {
		listener.Listen()
	}

	this.waiter.Done()
}
