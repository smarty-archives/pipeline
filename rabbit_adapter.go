package rabbit

import (
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

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
