package listeners

type CompositeListener struct {
	items []Listener
}

func NewCompositeListener(items ...Listener) *CompositeListener {
	filtered := make([]Listener, 0, len(items))

	for _, item := range items {
		if item != nil {
			filtered = append(filtered, item)
		}
	}

	return &CompositeListener{items: filtered}
}

func (this *CompositeListener) Listen() {
	length := len(this.items)
	if length == 0 {
		return
	}

	var (
		allButLast = this.items[0 : length-1]
		last       = this.items[length-1]
	)

	for _, item := range allButLast {
		go item.Listen()
	}

	last.Listen()
}

func (this *CompositeListener) Close() {
	for _, item := range this.items {
		if closer, ok := item.(ListenCloser); ok {
			closer.Close()
		}
	}
}
