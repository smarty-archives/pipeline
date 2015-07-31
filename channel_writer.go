package rabbit

type ChannelWriter struct {
	transactional bool
	controller    Controller
	channel       Channel
}

func newWriter(controller Controller, transactional bool) *ChannelWriter {
	return &ChannelWriter{
		controller:    controller,
		transactional: transactional,
	}
}

func (this *ChannelWriter) Write(message Dispatch) error {
	return nil
}

func (this *ChannelWriter) Commit() error {
	return nil
}

func (this *ChannelWriter) Close() {
	// this.controller.closeWriter(this)
	// within mutex mark as closed, but allow existing channel to proceed
	// another thread, while getting a channel
	// reference (within the mutex), checks to see
	// if we're closed
}
