package httpx

import "time"

type FakeWaiter struct {
	addCalls, doneCalls, waitCalls, counter int
	addCalled, waitCalled                   time.Time
}

func (this *FakeWaiter) Add(delta int) {
	this.addCalled = time.Now()
	this.addCalls++
	this.counter += delta
}

func (this *FakeWaiter) Done() {
	this.doneCalls++
	this.counter--
}

func (this *FakeWaiter) Wait() {
	this.waitCalls++
	this.waitCalled = time.Now()
}
