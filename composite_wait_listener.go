package listeners

import "sync"

type CompositeWaitListener struct {
	mutex  sync.Once
	waiter *sync.WaitGroup
	items  []Listener
}

func NewCompositeWaitShutdownListener(listeners ...Listener) *CompositeWaitListener {
	this := NewCompositeWaitListener()
	this.items = []Listener{NewShutdownListener(this.Close)}
	this.items = append(this.items, listeners...)
	return this
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

func (this *CompositeWaitListener) Close() {
	this.mutex.Do(this.close)
}

func (this *CompositeWaitListener) close() {
	for _, item := range this.items {
		if closer, ok := item.(ListenCloser); ok {
			closer.Close()
		}
	}
}
