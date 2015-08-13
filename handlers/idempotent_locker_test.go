package handlers

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type IdempotentLockerFixture struct {
	*gunit.Fixture

	inner  *FakeLocker
	locker *IdempotentLocker
}

func (this *IdempotentLockerFixture) Setup() {
	this.inner = &FakeLocker{}
	this.locker = NewIdempotentLocker(this.inner)
}

func (this *IdempotentLockerFixture) TestMultipleLocksOnlyCallsInnerOnce() {
	this.locker.Lock()
	this.locker.Lock()
	this.locker.Lock()
	this.locker.Lock()

	this.So(this.inner.locks, should.Equal, 1)
	this.So(this.inner.unlocks, should.Equal, 0)
}

func (this *IdempotentLockerFixture) TestUnlockWithoutLock() {
	this.locker.Unlock()
	this.locker.Unlock()
	this.locker.Unlock()

	this.So(this.inner.locks, should.Equal, 0)
	this.So(this.inner.unlocks, should.Equal, 0)
}

func (this *IdempotentLockerFixture) TestMultipleUnlocksOnlyCallsInnerOnce() {
	this.locker.Lock()
	this.locker.Unlock()
	this.locker.Unlock()
	this.locker.Unlock()

	this.So(this.inner.locks, should.Equal, 1)
	this.So(this.inner.unlocks, should.Equal, 1)
}

type FakeLocker struct{ locks, unlocks int }

func (this *FakeLocker) Lock()   { this.locks++ }
func (this *FakeLocker) Unlock() { this.unlocks++ }
