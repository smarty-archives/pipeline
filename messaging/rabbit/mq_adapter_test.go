package rabbit

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/pipeline/messaging"
	"github.com/streadway/amqp"
)

func TestRabbitAdapterFixture(t *testing.T) {
	gunit.Run(new(RabbitAdapterFixture), t)
}

type RabbitAdapterFixture struct {
	*gunit.Fixture
	now time.Time
}

func (this *RabbitAdapterFixture) Setup() {
	this.now = clock.UTCNow()
}

/////////////////////////////////////////////////////////////////////////////////

func (this *RabbitAdapterFixture) TestAMQPDeliveryConversion() {
	upstream := amqp.Delivery{
		AppId:           "1234",
		MessageId:       "5678",
		Type:            "message-type",
		ContentType:     "content-type",
		ContentEncoding: "content-encoding",
		Timestamp:       this.now,
		Body:            []byte{1, 2, 3, 4, 5, 6},
		DeliveryTag:     8675309,
	}

	this.So(fromAMQPDelivery(upstream, nil), should.Resemble, messaging.Delivery{
		SourceID:        1234,
		MessageID:       5678,
		MessageType:     "message-type",
		ContentType:     "content-type",
		ContentEncoding: "content-encoding",
		Timestamp:       this.now,
		Payload:         upstream.Body,
		Upstream:        upstream,
		Receipt:         DeliveryReceipt{channel: nil, deliveryTag: upstream.DeliveryTag},
	})
}

/////////////////////////////////////////////////////////////////////////////////

func (this *RabbitAdapterFixture) TestAMQPDispatchConversion() {
	dispatch := messaging.Dispatch{
		SourceID:        1234,
		MessageID:       5678,
		MessageType:     "message-type",
		ContentType:     "content-type",
		ContentEncoding: "content-encoding",
		Timestamp:       this.now.Add(-time.Second),
		Expiration:      time.Second,
		Durable:         true,
		Payload:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	this.So(toAMQPDispatch(dispatch, this.now), should.Resemble, amqp.Publishing{
		AppId:           "1234",
		MessageId:       "5678",
		Type:            "message-type",
		ContentType:     "content-type",
		ContentEncoding: "content-encoding",
		Timestamp:       this.now.Add(-time.Second),
		Expiration:      "1",
		DeliveryMode:    2,
		Body:            dispatch.Payload,
	})
}

/////////////////////////////////////////////////////////////////////////////////

func (this *RabbitAdapterFixture) TestAMQPDispatchTimestamp() {
	actual := toAMQPDispatch(messaging.Dispatch{}, this.now)
	expected := amqp.Publishing{
		AppId:        "0",
		MessageId:    "0",
		Timestamp:    this.now,
		DeliveryMode: amqp.Transient,
	}
	this.So(actual, should.Resemble, expected)
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
	this.assertExpiration(0, "")
	this.assertExpiration(time.Second, "1")
	this.assertExpiration(-time.Second, "0")
}
func (this *RabbitAdapterFixture) assertExpiration(expiration time.Duration, expected string) {
	this.So(computeExpiration(expiration), should.Equal, expected)
}

/////////////////////////////////////////////////////////////////////////////////

func (this *RabbitAdapterFixture) TestPersistenceComputation() {
	this.So(computePersistence(true), should.Equal, amqp.Persistent)
	this.So(computePersistence(false), should.Equal, amqp.Transient)
}

/////////////////////////////////////////////////////////////////////////////////
