//Building an e commerce system with 2 microservices: Order service and payment service
//oRDER SERVICE Creates order, payment service processe the payment

package saga_struct

type OrderCreatedEvent struct {
	OrderID string
	Amount  float64
}

type PaymentProcessedEvent struct {
	OrderID     string
	PaymentID   string
	PaymentDone bool
}

type CompensateOrderEvent struct {
	OrderID string
	Reason  string
}
