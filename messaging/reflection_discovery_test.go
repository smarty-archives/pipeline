package messaging

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type ReflectionDiscoveryFixture struct {
	*gunit.Fixture

	discovery *ReflectionDiscovery
}

func (this *ReflectionDiscoveryFixture) Setup() {
	this.discovery = NewReflectionDiscovery("prefix.")
}

////////////////////////////////////////////////////////////////////

func (this *ReflectionDiscoveryFixture) TestNilInstance() {
	result, err := this.discovery.Discover(nil)

	this.So(result, should.Equal, "")
	this.So(err, should.Equal, MessageTypeDiscoveryError)
}

////////////////////////////////////////////////////////////////////

func (this *ReflectionDiscoveryFixture) TestSimpleTypes() {
	this.assertDiscovery(uint64(0), "prefix.uint64")
	this.assertDiscovery(0, "prefix.int")
	this.assertDiscovery(true, "prefix.bool")
	this.assertDiscovery(SampleMessage{}, "prefix.samplemessage")
	this.assertDiscovery(&SampleMessage{}, "*prefix.samplemessage")
}
func (this *ReflectionDiscoveryFixture) assertDiscovery(instance interface{}, expected string) {
	result, err := this.discovery.Discover(instance)

	this.So(result, should.Equal, expected)
	this.So(err, should.BeNil)
}

////////////////////////////////////////////////////////////////////

func (this *ReflectionDiscoveryFixture) TestAnonymousStructs() {
	anonymous := struct{ Value1, Value2, Value3 string }{}
	result, err := this.discovery.Discover(anonymous)
	this.So(result, should.Equal, "")
	this.So(err, should.Equal, MessageTypeDiscoveryError)
}

////////////////////////////////////////////////////////////////////

func (this *ReflectionDiscoveryFixture) TestEmptyStructs() {
	result, err := this.discovery.Discover(struct{}{})
	this.So(result, should.Equal, "")
	this.So(err, should.Equal, MessageTypeDiscoveryError)
}

////////////////////////////////////////////////////////////////////

type SampleMessage struct{}
