package postnord

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/cadence"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

const (
	WorkflowName = "OrderProcessingWorkflow"
)

// RegisterWorkflow registers the OrderProcessingWorkflow.
func RegisterWorkflow(w worker.Worker) {
	w.RegisterWorkflowWithOptions(OrderProcessingWorkflow, workflow.RegisterOptions{Name: WorkflowName})

	// Register your activities here
	w.RegisterActivityWithOptions(validatePayment, activity.RegisterOptions{Name: "validatePayment"})
	w.RegisterActivityWithOptions(shipPackage, activity.RegisterOptions{Name: "shipPackage"})
}

// Order represents an order with basic details like the ID, customer name, and order amount.
type Order struct {
	ID       string  `json:"id"`
	Customer string  `json:"customer"`
	Amount   float64 `json:"amount"`
	Address  string  `json:"address"`
}

// OrderProcessingWorkflow processes an order through several steps:
// 1. It first validates the payment for the order.
// 2. Then, it proceeds to ship the package.
// 3. Finally, it returns a result indicating success or failure based on the payment and shipping status.
func OrderProcessingWorkflow(ctx workflow.Context, order Order) (string, error) {
	retryPolicy := &cadence.RetryPolicy{
		InitialInterval:    1 * time.Second,  // Start with 1 second.
		BackoffCoefficient: 2.0,              // Exponential backoff.
		MaximumInterval:    10 * time.Second, // Max retry interval.
		MaximumAttempts:    3,                // Retry up to 3 times.
	}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		RetryPolicy:            retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Starting OrderProcessingWorkflow", zap.String("orderID", order.ID), zap.String("customer", order.Customer))
	fmt.Print("System print: Starting OrderProcessingWorkflow", zap.String("orderID", order.ID), zap.String("customer", order.Customer))

	// Step 1: Validate the payment.
	// The payment validation step checks if the payment for the order is valid.
	// In this example, we simulate the payment validation by calling the `validatePayment` activity.
	// If validation fails, the workflow stops early and returns an appropriate error.
	var paymentValidationResult string
	err := workflow.ExecuteActivity(ctx, validatePayment, order).Get(ctx, &paymentValidationResult)
	if err != nil {
		return "", fmt.Errorf("validate payment for order: %v", err)
	}

	workflow.Sleep(ctx, 2*time.Minute)
	fmt.Print("System print2: Starting OrderProcessingWorkflow", zap.String("orderID", order.ID), zap.String("customer", order.Customer))

	// Step 2: Ship the package
	// Once the payment is validated, we proceed to ship the package.
	// The ship the package activity is called to simulate the shipping process.
	// If shipping fails or encounters an error, the workflow returns an error.

	// Step 3: Return a success message
	// If both payment validation and shipping were successful, we return a success message indicating the order was processed.
	return fmt.Sprintf("Order %s processed successfully for customer %s.", order.ID, order.Customer), nil
}

// Add an activity here that validates a payment.
// The validation fails if the order amount is greater than 25 (for example, due to payment policy restrictions).
func validatePayment(ctx context.Context, order Order) (string, error) {
	if order.Amount > 25 {
		return "", fmt.Errorf("payment validation failed: order amount exceeds limit")
	}

	info := activity.GetInfo(ctx)
	if info.Attempt < 2 {
		activity.GetLogger(ctx).Info("Temporary failure in payment processing", zap.Int32("attempt", info.Attempt))
		return "", fmt.Errorf("temporary issue, please retry")
	}

	activity.GetLogger(ctx).Info("Payment processed successfully")

	return "Payment successfull!", nil
}

// Add another activity that ships the package.
// This activity checks if the shipping address is provided. The shipping fails if the address is empty.
func shipPackage(ctx context.Context, order Order) (string, error) {

	if len(order.Customer) == 0 {
		return "", fmt.Errorf("shipping failed: customer (shipping address) is missing")
	}
	return "Shipment successfull!", nil
}
