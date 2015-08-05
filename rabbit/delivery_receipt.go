package rabbit

type DeliveryReceipt struct {
	channel     Acknowledger
	deliveryTag uint64
}

func newReceipt(channel Acknowledger, deliveryTag uint64) DeliveryReceipt {
	return DeliveryReceipt{
		channel:     channel,
		deliveryTag: deliveryTag,
	}
}
