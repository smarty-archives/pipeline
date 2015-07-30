package rabbit

type DeliveryReceipt struct {
	channel     Channel
	deliveryTag uint64
}

func newReceipt(channel Channel, deliveryTag uint64) DeliveryReceipt {
	return DeliveryReceipt{
		channel:     channel,
		deliveryTag: deliveryTag,
	}
}
