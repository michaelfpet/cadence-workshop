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

package client

import (
	"context"
	"fmt"

	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/fx"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/grpc"
)

const (
	hostAddress    = "127.0.0.1:7833"
	clientName     = "cadence-workshop-worker"
	cadenceService = "cadence-frontend"
)

type Params struct {
	fx.In

	Lc fx.Lifecycle
}

func New(params Params) workflowserviceclient.Interface {
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: clientName,
		Outbounds: yarpc.Outbounds{
			cadenceService: {Unary: grpc.NewTransport().NewSingleOutbound(hostAddress)},
		},
	})

	setUpLcHooks(params, dispatcher)
	clientConfig := dispatcher.ClientConfig(cadenceService)

	return workflowserviceclient.New(clientConfig)
}

var Module = fx.Provide(New)

func setUpLcHooks(params Params, dispatcher *yarpc.Dispatcher) {
	params.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				err := dispatcher.Start()
				if err != nil {
					return fmt.Errorf("start dispatcher: %w", err)
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				err := dispatcher.Stop()
				if err != nil {
					return fmt.Errorf("stop dispatcher: %w", err)
				}
				return nil
			},
		})
}
