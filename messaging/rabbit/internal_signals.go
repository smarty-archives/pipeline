package rabbit

type shutdownRequested struct{}
type subscriptionClosed struct {
	DeliveryCount     uint64
	LatestDeliveryTag uint64
	LatestConsumer    interface{}
}
type acknowledgementCompleted struct{ Acknowledgements uint64 }
