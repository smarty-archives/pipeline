package rabbit

type shutdownRequested struct{}
type subscriptionClosed struct{ DeliveryCount uint64 }
type acknowledgementCompleted struct{ Acknowledgements uint64 }
