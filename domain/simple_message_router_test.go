package domain

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestSimpleMessageRouterFixture(t *testing.T) {
	gunit.Run(new(SimpleMessageRouterFixture), t)
}

type SimpleMessageRouterFixture struct {
	*gunit.Fixture

	messageRouter *SimpleMessageRouter
	aggregate1    *FakeAggregate
	aggregate2    *FakeAggregate
}

func (this *SimpleMessageRouterFixture) Setup() {
	this.aggregate1 = NewFakeAggregate("aggregate1")
	this.aggregate2 = NewFakeAggregate("aggregate2")
	this.messageRouter = NewSimpleMessageRouter(
		[]Handler{this.aggregate1, this.aggregate2},
		[]Applicator{this.aggregate1, this.aggregate2})
}

func (this *SimpleMessageRouterFixture) TestEachDocumentIsPresentedWithTheEvent() {
	this.messageRouter.Apply("Message")

	for _, document := range []*FakeAggregate{this.aggregate1, this.aggregate2} {
		this.So(document.applied, should.Resemble, []interface{}{"Message"})
	}
}

func (this *SimpleMessageRouterFixture) TestWhenNoDocumentsApplyTheMessageTheMessageRouterReturnsFalse() {
	result := this.messageRouter.Apply("Message")
	this.So(result, should.BeFalse)
}

func (this *SimpleMessageRouterFixture) TestWhenAnyDocumentAppliesTheMessageTheRouterReturnsTrue() {
	this.aggregate2.canReceive = true
	result := this.messageRouter.Apply("Message")
	this.So(result, should.BeTrue)
}

func (this *SimpleMessageRouterFixture) TestCommandsAreHandledAndResultingEventsAreAppliedAndReturned() {
	this.aggregate1.canReceive = true
	this.aggregate2.canReceive = true

	actualEvents := this.messageRouter.Handle("Command")

	expectedEvents := []interface{}{
		"aggregate1" + "Command" + "event1",
		"aggregate1" + "Command" + "event2",
		"aggregate2" + "Command" + "event1",
		"aggregate2" + "Command" + "event2",
	}
	this.So(actualEvents, should.Resemble, expectedEvents)
	this.So(this.aggregate1.applied, should.Resemble, expectedEvents)
	this.So(this.aggregate2.applied, should.Resemble, expectedEvents)
}

///////////////////////////////////////////////////////////////////////////////

func (this *SimpleMessageRouterFixture) TestHandlersCanBeAdded() {
	handler := NewFakeAggregate("")
	this.messageRouter.Add(handler)

	this.messageRouter.Handle("1")

	this.So(handler.handled, should.Resemble, []interface{}{"1"})
}

///////////////////////////////////////////////////////////////////////////////

func (this *SimpleMessageRouterFixture) TestDocumentsCanBeAdded() {
	document := NewFakeAggregate("")
	this.messageRouter.Add(document)

	this.messageRouter.Apply("2")

	this.So(document.applied, should.Resemble, []interface{}{"2"})
}

///////////////////////////////////////////////////////////////////////////////

func (this *SimpleMessageRouterFixture) TestCompositeItemCanBeAdded() {
	document := NewFakeAggregate("")
	this.messageRouter = NewSimpleMessageRouter(nil, nil)

	added := this.messageRouter.Add(document)
	this.messageRouter.Handle("1")
	this.messageRouter.Apply("abc")

	this.So(added, should.BeTrue)
	this.So(document.handled, should.Resemble, []interface{}{"1"})
	this.So(document.applied, should.Resemble, []interface{}{"1event1", "1event2", "abc"})
}

///////////////////////////////////////////////////////////////////////////////

func (this *SimpleMessageRouterFixture) TestHandleMultipleTimesGivesBackCorrectEvents() {
	handler := NewFakeAggregate("")
	this.messageRouter.Add(handler)

	events1 := this.messageRouter.Handle("1")
	events2 := this.messageRouter.Handle("1")

	this.So(events1, should.Resemble, events2)
}

///////////////////////////////////////////////////////////////////////////////

type FakeAggregate struct {
	id         string
	handled    []interface{}
	applied    []interface{}
	canReceive bool
}

func NewFakeAggregate(id string) *FakeAggregate {
	return &FakeAggregate{id: id, applied: []interface{}{}}
}

func (this *FakeAggregate) Handle(command interface{}) []interface{} {
	this.handled = append(this.handled, command)
	return []interface{}{
		this.id + command.(string) + "event1",
		this.id + command.(string) + "event2",
	}
}

func (this *FakeAggregate) Apply(event interface{}) bool {
	this.applied = append(this.applied, event)
	return this.canReceive
}
