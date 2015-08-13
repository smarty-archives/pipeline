package domain

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type SimpleMessageRouterFixture struct {
	*gunit.Fixture

	messageRouter *SimpleMessageRouter
	doc1          *FakeDocument
	doc2          *FakeDocument
}

func (this *SimpleMessageRouterFixture) Setup() {
	this.doc1 = NewMockApplicator("doc1")
	this.doc2 = NewMockApplicator("doc2")
	this.messageRouter = NewSimpleMessageRouter(
		[]Handler{this.doc1, this.doc2},
		[]Applicator{this.doc1, this.doc2})
}

func (this *SimpleMessageRouterFixture) TestEachDocumentIsPresentedWithTheEvent() {
	this.messageRouter.Apply("Message")

	for _, document := range []*FakeDocument{this.doc1, this.doc2} {
		this.So(document.received, should.Resemble, []interface{}{"Message"})
	}
}

func (this *SimpleMessageRouterFixture) TestWhenNoDocumentsApplyTheMessageTheMessageRouterReturnsFalse() {
	result := this.messageRouter.Apply("Message")
	this.So(result, should.BeFalse)
}

func (this *SimpleMessageRouterFixture) TestWhenAnyDocumentAppliesTheMessageTheRouterReturnsTrue() {
	this.doc2.canReceive = true
	result := this.messageRouter.Apply("Message")
	this.So(result, should.BeTrue)
}

func (this *SimpleMessageRouterFixture) TestCommandsAreHandledAndResultingEventsAreAppliedAndReturned() {
	this.doc1.canReceive = true
	this.doc2.canReceive = true

	actualEvents := this.messageRouter.Handle("Command")

	expectedEvents := []interface{}{
		"doc1" + "Command" + "event1",
		"doc1" + "Command" + "event2",
		"doc2" + "Command" + "event1",
		"doc2" + "Command" + "event2",
	}
	this.So(actualEvents, should.Resemble, expectedEvents)
	this.So(this.doc1.received, should.Resemble, expectedEvents)
	this.So(this.doc2.received, should.Resemble, expectedEvents)
}

///////////////////////////////////////////////////////////////////////////////

type FakeDocument struct {
	id         string
	received   []interface{}
	canReceive bool
}

func NewMockApplicator(id string) *FakeDocument {
	return &FakeDocument{id: id, received: []interface{}{}}
}

func (this *FakeDocument) Handle(command interface{}) []interface{} {
	return []interface{}{
		this.id + command.(string) + "event1",
		this.id + command.(string) + "event2",
	}
}

func (this *FakeDocument) Apply(event interface{}) bool {
	this.received = append(this.received, event)
	return this.canReceive
}
