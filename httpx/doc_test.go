package httpx

//go:generate go install github.com/smartystreets/gunit/gunit
//go:generate gunit

type FakeWaiter struct{ addCalls, doneCalls, counter int }

func (this *FakeWaiter) Add(delta int) {
	this.addCalls++
	this.counter += delta
}

func (this *FakeWaiter) Done() {
	this.doneCalls++
	this.counter--
}
