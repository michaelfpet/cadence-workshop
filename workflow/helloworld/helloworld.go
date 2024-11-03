// The MIT License (MIT)

// Copyright (c) 2017-2020 Uber Technologies Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package helloworld

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

const (
	WorkflowName = "HelloWorld"
	ActivityName = "HelloWorldActivity"
)

// RegisterWorkflow registers the HelloWorldWorkflow.
func RegisterWorkflow(w worker.Worker) {
	w.RegisterWorkflowWithOptions(HelloWorldWorkflow, workflow.RegisterOptions{Name: WorkflowName})
	w.RegisterActivityWithOptions(HelloWorldActivity, activity.RegisterOptions{
		Name:                ActivityName,
		EnableAutoHeartbeat: true,
	})
}

type helloWorldInput struct {
	Message string `json:"message"`
}

// HelloWorldWorkflow greets the caller.
func HelloWorldWorkflow(ctx workflow.Context, input helloWorldInput) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorldWorkflow started")

	var greetingMsg string
	err := workflow.ExecuteActivity(ctx, HelloWorldActivity, input).Get(ctx, &greetingMsg)
	if err != nil {
		logger.Error("HelloWorldActivity failed", zap.Error(err))
		return "", err
	}

	logger.Info("Workflow result", zap.String("greeting", greetingMsg))
	return greetingMsg, nil
}

// HelloWorldActivity constructs the greeting message from input.
func HelloWorldActivity(ctx context.Context, input helloWorldInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("HelloWorldActivity started")
	return fmt.Sprintf("Hello, %s!", input.Message), nil
}
