package saga_service

import (
	saga_struct "distributed_trans_custom/saga-pattern/struct_service"
	"fmt"
)

func CreateOrder(orderID string, amount float64) saga_struct.OrderCreatedEvent {
	return saga_struct.OrderCreatedEvent{OrderID: orderID, Amount: amount}
}

func ProcessPayment(orderID string, amount float64) saga_struct.PaymentProcessedEvent {
	return saga_struct.PaymentProcessedEvent{OrderID: orderID, PaymentDone: true}
}

func RollBackOrder(orderID string, reason string) {
	fmt.Println("Compensating order: ", orderID, "Reason", reason)
}

func HandlerOrderSaga(orderID string, amount float64) {
	orderCreated := CreateOrder(orderID, amount)
	if orderCreated.OrderID != "" {
		paymentProcessed := ProcessPayment(orderCreated.OrderID, orderCreated.Amount)
		if !paymentProcessed.PaymentDone {
			RollBackOrder(orderID, "Payment Failed")
		}
	}
}
