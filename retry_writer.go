package messenger

type RetryWriter struct {
	inner Writer
	max   uint64
	sleep func()
}

func NewRetryWriter(inner Writer, max uint64, sleep func()) *RetryWriter {
	if max == 0 {
		max = 0xFFFFFFFFFFFFFFFF
	}

	return &RetryWriter{
		inner: inner,
		max:   max,
		sleep: sleep,
	}
}
func (this *RetryWriter) Write(message Dispatch) (err error) {
	for i := uint64(0); i < this.max; i++ {
		if err = this.inner.Write(message); err == nil {
			break
		} else if err == WriterClosedError {
			break
		} else {
			this.sleep()
		}
	}

	return err
}

func (this *RetryWriter) Close() {
	this.inner.Close()
}
