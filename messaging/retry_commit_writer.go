package messaging

type RetryCommitWriter struct {
	inner   CommitWriter
	max     uint64
	success func(uint64)
	sleep   func(uint64)
	buffer  []Dispatch
}

func NewRetryCommitWriter(inner CommitWriter, max uint64, success func(uint64), sleep func(uint64)) *RetryCommitWriter {
	if max == 0 {
		max = 0xFFFFFFFFFFFFFFFF
	}

	if success == nil {
		success = func(uint64) {}
	}

	return &RetryCommitWriter{
		inner:   inner,
		max:     max,
		success: success,
		sleep:   sleep,
	}
}

func (this *RetryCommitWriter) Write(message Dispatch) error {
	this.buffer = append(this.buffer, message)
	return nil
}

func (this *RetryCommitWriter) Commit() (err error) {
	for i := uint64(0); i < this.max; i++ {
		if err = this.try(); err == nil {
			this.success(i)
			this.buffer = this.buffer[0:0]
			return nil
		} else if err == WriterClosedError {
			return err
		}

		this.sleep(i)
	}

	return err
}
func (this *RetryCommitWriter) try() error {
	for _, item := range this.buffer {
		if this.inner.Write(item) == WriterClosedError {
			return WriterClosedError
		}
	}

	return this.inner.Commit()
}

func (this *RetryCommitWriter) Close() {
	this.inner.Close()
}
