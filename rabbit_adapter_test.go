package rabbit

import (
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/streadway/amqp"
)

type RabbitAdapterFixture struct {
	*gunit.Fixture
	now time.Time
}

func (this *RabbitAdapterFixture) Setup() {
	this.now = time.Now()
}

/////////////////////////////////////////////////////////////////////////////////

func (this *RabbitAdapterFixture) TestParsingNumericString() {
	this.assertParsedValue("1", 1)
	this.assertParsedValue("", 0)
	this.assertParsedValue("-1", 0)
	this.assertParsedValue("-2", 0)
	this.assertParsedValue("18446744073709551615", 18446744073709551615)
}
func (this *RabbitAdapterFixture) assertParsedValue(value string, expected uint64) {
	this.So(parseUint64(value), should.Equal, expected)
}

/////////////////////////////////////////////////////////////////////////////////

func (this *RabbitAdapterFixture) TestExpirationComputation() {
	this.assertExpiration(time.Time{}, "")
	this.assertExpiration(this.now.Add(time.Second), "1")
	this.assertExpiration(this.now.Add(-time.Second), "0")
}
func (this *RabbitAdapterFixture) assertExpiration(expiration time.Time, expected string) {
	this.So(computeExpiration(this.now, expiration), should.Equal, expected)
}

/////////////////////////////////////////////////////////////////////////////////

func (this *RabbitAdapterFixture) TestPersistenceComputation() {
	this.So(computePersistence(true), should.Equal, amqp.Persistent)
	this.So(computePersistence(false), should.Equal, amqp.Transient)
}
