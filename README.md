# Getting started

You will need three terminals to run the code in this repository.
1. One for running the Cadence server via Docker Compose.
2. One for running the Cadence CLI.
3. One for running the Cadence worker contained in this repository.

## Running the Cadence server

In one terminal, run the following command from the repository root:
```bash
docker-compose up
```

You will see a lot of output as the Cadence server starts up. 
The cadence web UI will be available at http://localhost:8088.

## Running the Cadence CLI

First lets setup an alias for the Cadence CLI.
### Alias for Cadence CLI
An easy way to run the Cadence CLI is to use Docker. We will
create an alias for the command to make it easier to run. If you do not 
want to create an alias, you can run the command directly.

After setting up the alias you can run the Cadence CLI by running `cadence` in a terminal.

#### Mac:
Run the following command in a terminal:
```bash
alias cadence="docker run -i -t ubercadence/cli:master --address host.docker.internal:7933"
```

#### Linux:
Run the following command in a terminal:
```bash
alias cadence="docker run -i -t --network=host ubercadence/cli:master --address 127.0.0.1:7933"
```

#### Windows:
Create a file called `cadence.bat` with the following content:
```bash
docker run -i -t ubercadence/cli:master --address host.docker.internal:7933 %*
```

### Creating a domain for the worker
Before we can start the worker, we need to create a domain for it.
Run the following command in a terminal:
```bash
cadence --domain cadence-workshop domain register
```

You should now be able to see the domain in the Cadence web UI here:
http://localhost:8088/domains/cadence-workshop

## Running the Cadence worker
To run the cadence worker, simply run the following from the repository root:
```bash
go run main.go
```

## Starting a workflow
To start the hello world workflow, run the following command:
```bash
cadence --domain cadence-workshop wf start --tasklist tasklist --execution_timeout 10 --workflow_type HelloWorld --input '{"message":"Cadence"}'
```

The workflow should start, and you should see the workflow in the web UI:
http://localhost:8088/domains/cadence-workshop

### Task 1: Postnord OrderProcessingWorkflow
Implement the validatePayment() and shipProduct() activities for the OrderProcessingWorkflow. <br />
The validatePayment() activity should check if the order amount is greater than 25 and return an error if it is, otherwise return true for valid payments. <br />
The shipProduct() activity should verify that the order.Customer (representing the shipping address) is not empty; <br />
if it is, return an error indicating that shipping failed, otherwise simulate shipping and return true for successful shipment. <br />
To start the OrderProcessingWorkflow, run the following command:
```bash
cadence --domain cadence-workshop wf start --tasklist tasklist --execution_timeout 10 --workflow_type OrderProcessingWorkflow --input '{"id":"Order123", "customer": "Cadence", "amount": 20, "address": "Uber office"}'
```

### Task 2: Simulate Payment Failure in validatePayment() activity
First, modify the validatePayment to simulate a failure. This can be done by causing the activity to fail intermittently, such as by returning an error during the first few attempts, which will simulate a temporary issue like a network failure.
<br />
<br />
Use the code below to simulate the failure. <br />
info.Attempt value from the activity.GetInfo(ctx) function tracks the current attempt count. 
```go
func validatePayment(ctx context.Context, order Order) (string, error) {
	// Simulate a failure
	info := activity.GetInfo(ctx)
	if info.Attempt < 3 {
		activity.GetLogger(ctx).Info("Temporary failure in payment processing")
		return 0, fmt.Errorf("temporary issue, please retry")
	}
	
	activity.GetLogger(ctx).Info("Payment processed successfully")
	return "Payment successful", nil
}
```

Ensure that the workflow fails in the first runs and returns an error.

### Task 3: Add Retry Policy to validatePayment() activity
Now, implement a retry policy to handle the temporary failures simulated in Task 2. <br /> 
The retry policy should retry the validatePayment activity up to 3 times with exponential backoff. <br />
Use the paymentRetryPolicy configuration below:
Read more about the activity and workflow retries: https://cadenceworkflow.io/docs/go-client/retries/
```go
// Retry policy configuration: exponential backoff with a maximum of 3 retries.
var paymentRetryPolicy = &activity.RetryPolicy{
	InitialInterval:    1 * time.Second,    // Start with 1 second.
	BackoffCoefficient: 2.0,                // Exponential backoff.
	MaximumInterval:    10 * time.Second,   // Max retry interval.
	MaximumAttempts:    3,                  // Retry up to 3 times.
}
```
retry policy configuration can be applied to an activity like this:
```go
// Configure activity options with retry policy.
activityOptions := workflow.ActivityOptions{	
	RetryPolicy: paymentRetryPolicy, // Attach retry policy.
}

// Set the activity options to the context.
activityCtx = workflow.WithActivityOptions(ctx, activityOptions)

// Start the activity
workflow.ExecuteActivity(activityCtx, activities.validatePayment)
```
- workflow.ActivityOptions{}: This is a struct in Cadence that allows you to configure options for executing an activity within a workflow.
- workflow.WithActivityOptions(ctx, activityOptions): This is a function provided by Cadence that applies the given activity options (such as the retry policy) to the workflow context (ctx).

Ensure that now the workflow retries validatePayment() activity.
