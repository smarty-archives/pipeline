package rabbit

type shutdownRequested struct{}
type subscriptionClosed struct{ Deliveries uint64 }
type acknowledgementCompleted struct{ Acknowledgements uint64 }
