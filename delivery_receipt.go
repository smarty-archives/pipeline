package rabbit

type DeliveryReceipt struct {
	channel     Channel
	deliveryTag uint64
}
