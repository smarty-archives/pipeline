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
