package rabbit

type shutdownRequested struct{}
type subscriptionClosed struct{ FinalReceipt interface{} }
type acknowledgementCompleted struct{}
