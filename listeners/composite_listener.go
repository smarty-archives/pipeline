package listeners

type CompositeListener struct {
	listeners []Listener
}

func NewCompositeListener(listeners ...Listener) Listener {
	filtered := make([]Listener, 0, len(listeners))
	for _, listener := range listeners {
		if listener != nil {
			filtered = append(filtered, listener)
		}
	}

	return &CompositeListener{listeners: filtered}
}

func (this *CompositeListener) Listen() {
	length := len(this.listeners)
	if length == 0 {
		return
	}

	var (
		allButLast = this.listeners[0 : length-1]
		last       = this.listeners[length-1]
	)

	for _, item := range allButLast {
		go item.Listen()
	}

	last.Listen()
}
