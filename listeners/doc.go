package listeners

// Listener generally describes a struct that listens for items from
// one or more input channels (received in it's constructor) and
// optionally sends other items on one or more output channels.
type Listener interface {
	// Listen is generally a long-running method, and should be executed within
	// its own goroutine so that other Listeners may run concurrently.
	Listen()
}

type WaitGroup interface {
	Add(delta int)
	Done()
}
