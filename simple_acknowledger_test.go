package rabbit

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type SimpleAcknowledgerFixture struct {
	*gunit.Fixture

	channel      *FakeAcknowledgmentChannel
	acknowledger *BatchAcknowledger
	control      chan interface{}
	input        chan interface{}
}

func (this *SimpleAcknowledgerFixture) Setup() {
	this.channel = &FakeAcknowledgmentChannel{}
	this.control = make(chan interface{}, 16)
	this.input = make(chan interface{}, 16)
	this.acknowledger = newAcknowledger(this.control, this.input)
	go this.acknowledger.Listen()
}

////////////////////////////////////////////////////////////////

func (this *SimpleAcknowledgerFixture) TestItemIsAcknowledged() {
	this.input <- newReceipt(this.channel, 5678)

	close(this.input)
	time.Sleep(time.Millisecond)

	this.So(this.channel.callsMulti, should.Equal, 1)
	this.So(this.channel.latestMulti, should.Equal, 5678)
}

////////////////////////////////////////////////////////////////

func (this *SimpleAcknowledgerFixture) TestOnlyLastItemIsAcknowledged() {
	this.input <- newReceipt(this.channel, 5678)
	this.input <- newReceipt(this.channel, 5679)

	close(this.input)
	time.Sleep(time.Millisecond)

	this.So(this.channel.callsMulti, should.Equal, 1)
	this.So(this.channel.latestMulti, should.Equal, 5679)
}

////////////////////////////////////////////////////////////////

func (this *SimpleAcknowledgerFixture) TestControlChannelReceivesCompletionNotice() {
	this.input <- newReceipt(this.channel, 1)
	this.input <- newReceipt(this.channel, 2)
	this.input <- newReceipt(this.channel, 3)

	close(this.input)

	this.So((<-this.control).(acknowledgementCompleted).Acknowledgements, should.Equal, 3)
}

////////////////////////////////////////////////////////////////

func (this *SimpleAcknowledgerFixture) TestFinalReceiptIsAlwaysCalled() {
	this.input <- newReceipt(this.channel, 1)
	this.input <- newReceipt(this.channel, 2)
	this.input <- newReceipt(this.channel, 3)
	this.input <- 0 // junk to be ignored

	close(this.input)

	this.So((<-this.control).(acknowledgementCompleted).Acknowledgements, should.Equal, 3)
	this.So(this.channel.callsMulti, should.Equal, 1)
	this.So(this.channel.latestMulti, should.Equal, 3)
}

////////////////////////////////////////////////////////////////

func (this *SimpleAcknowledgerFixture) TestLoopExitsAfterAllDeliveriesReceived() {
	this.input <- newReceipt(this.channel, 1)
	this.input <- newReceipt(this.channel, 2)
	this.input <- newReceipt(this.channel, 3)
	this.input <- subscriptionClosed{Deliveries: 1}
	this.input <- subscriptionClosed{Deliveries: 2}

	this.So((<-this.control).(acknowledgementCompleted).Acknowledgements, should.Equal, 3)
	this.So(this.channel.callsMulti, should.Equal, 1)
	this.So(this.channel.latestMulti, should.Equal, 3)
}

////////////////////////////////////////////////////////////////

func (this *SimpleAcknowledgerFixture) TestLoopExitsAfterAllDeliveriesReceived2() {
	this.input <- subscriptionClosed{Deliveries: 1}
	this.input <- newReceipt(this.channel, 1)
	this.input <- subscriptionClosed{Deliveries: 2}
	this.input <- newReceipt(this.channel, 2)
	this.input <- newReceipt(this.channel, 3)

	this.So((<-this.control).(acknowledgementCompleted).Acknowledgements, should.Equal, 3)
	this.So(this.channel.callsMulti, should.Equal, 1)
	this.So(this.channel.latestMulti, should.Equal, 3)
}

////////////////////////////////////////////////////////////////

type FakeAcknowledgmentChannel struct {
	calls       uint64
	callsMulti  uint64
	callsSingle uint64

	latest       uint64
	latestMulti  uint64
	latestSingle uint64

	tags       []uint64
	tagsMulti  []uint64
	tagsSingle []uint64
}

func (this *FakeAcknowledgmentChannel) AcknowledgeSingleMessage(value uint64) error {
	this.callsSingle++
	this.latestSingle = value
	this.tagsSingle = append(this.tagsSingle, value)

	return this.acknowledgeMessage(value)
}
func (this *FakeAcknowledgmentChannel) AcknowledgeMultipleMessages(value uint64) error {
	this.callsMulti++
	this.latestMulti = value
	this.tagsMulti = append(this.tagsMulti, value)

	return this.acknowledgeMessage(value)
}
func (this *FakeAcknowledgmentChannel) acknowledgeMessage(value uint64) error {
	this.calls++
	this.latest = value
	this.tags = append(this.tags, value)

	return nil
}
