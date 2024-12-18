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

package worker

import (
	"context"
	"fmt"

	"github.com/uber-go/tally"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	domain = "cadence-workshop"

	taskListName = "tasklist"
)

type Params struct {
	fx.In

	Lc            fx.Lifecycle
	Logger        *zap.Logger
	Metrics       tally.Scope
	CadenceClient workflowserviceclient.Interface
}

type Result struct {
	fx.Out

	Worker worker.Worker
}

func New(params Params) Result {
	workerOptions := worker.Options{
		Logger:       params.Logger,
		MetricsScope: params.Metrics,
	}

	w := worker.New(params.CadenceClient, domain, taskListName, workerOptions)

	setUpLcHooks(params.Lc, w)

	return Result{
		Worker: w,
	}
}

var Module = fx.Provide(New)

func setUpLcHooks(lc fx.Lifecycle, w worker.Worker) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := w.Start()
			if err != nil {
				return fmt.Errorf("start worker: %w", err)
			}
			return nil
		},
		OnStop: func(context.Context) error {
			w.Stop()
			return nil
		},
	})
}
