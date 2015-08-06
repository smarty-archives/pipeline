package rabbit

import (
	"strconv"
	"time"

	"github.com/smartystreets/go-messenger"
	"github.com/streadway/amqp"
)

func fromAMQPDelivery(delivery amqp.Delivery, channel Acknowledger) messenger.Delivery {
	return messenger.Delivery{
		SourceID:    parseUint64(delivery.AppId),
		MessageID:   parseUint64(delivery.MessageId),
		MessageType: delivery.Type,
		Encoding:    delivery.ContentEncoding,
		Payload:     delivery.Body,
		Upstream:    delivery,
		Receipt:     newReceipt(channel, delivery.DeliveryTag),
	}
}
func parseUint64(value string) uint64 {
	parsed, _ := strconv.ParseUint(value, 10, 64)
	return parsed
}

func toAMQPDispatch(dispatch messenger.Dispatch, now time.Time) amqp.Publishing {
	return amqp.Publishing{
		AppId:           strconv.FormatUint(dispatch.SourceID, base10),
		MessageId:       strconv.FormatUint(dispatch.MessageID, base10),
		Type:            dispatch.MessageType,
		ContentEncoding: dispatch.Encoding,
		Timestamp:       now,
		Expiration:      computeExpiration(now, dispatch.Expiration),
		DeliveryMode:    computePersistence(dispatch.Durable),
		Body:            dispatch.Payload,
	}
}
func computeExpiration(now, expiration time.Time) string {
	if expiration == noExpiration {
		return ""
	} else if seconds := expiration.Sub(now).Seconds(); seconds <= 0 {
		return "0"
	} else {
		return strconv.FormatUint(uint64(seconds), base10)
	}
}
func computePersistence(durable bool) uint8 {
	if durable {
		return amqp.Persistent
	}

	return amqp.Transient
}

var noExpiration = time.Time{}

const base10 = 10
