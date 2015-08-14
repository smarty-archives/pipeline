package handlers

import "sync"

/// This implementation is not safe to be used by multiple goroutines
/// it is meant exist in and by owned by a single goroutine.
/// Its purpose is to prevent the underlying locker from being locked
/// multiple times.
type IdempotentLocker struct {
	locked bool
	locker sync.Locker
}

func NewIdempotentLocker(locker sync.Locker) sync.Locker {
	if locker == nil {
		return NoopLocker{}
	}
	return &IdempotentLocker{locked: false, locker: locker}
}

func (this *IdempotentLocker) Lock() {
	if !this.locked {
		this.locked = true
		this.locker.Lock()
	}
}

func (this *IdempotentLocker) Unlock() {
	if this.locked {
		this.locked = false
		this.locker.Unlock()
	}
}
