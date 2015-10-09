package listeners

type SerialListener struct {
	items []Listener
}

func NewSerialListener(items ...Listener) *SerialListener {
	return &SerialListener{items: items}
}

func (this *SerialListener) Listen() {
	for _, item := range this.items {
		if item != nil {
			item.Listen()
		}
	}
}

func (this *SerialListener) Close() {
	for _, item := range this.items {
		if closer, ok := item.(ListenCloser); ok {
			closer.Close()
		}
	}
}
