package handlers

type NoopLocker struct{}

func (this NoopLocker) Lock()   {}
func (this NoopLocker) Unlock() {}
