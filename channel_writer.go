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
	// WARNING: close can (and will) be called from other threads
}
