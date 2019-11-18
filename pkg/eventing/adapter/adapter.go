/*
Copyright 2019 The Knative Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This package contains proposed extensions to "knative.dev/eventing/pkg/adapter"
// For now we include them with the sample source.
package adapter

import (
	"context"

	"github.com/cloudevents/sdk-go/pkg/binding"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"go.uber.org/zap"
	"knative.dev/eventing/pkg/adapter"
	"knative.dev/pkg/logging"
)

type EnvConfigAccessor = adapter.EnvConfigAccessor
type EnvConfig = adapter.EnvConfig
type Adapter = adapter.Adapter

var Main = adapter.Main

// WithCancelChannel is like context.WithCancel but also arranges that closing stopCh
// will call CancelFunc.
//
// Note the caller must *always* call the CancelFunc (directly or by closing stopCh)
// otherwise resources may be leaked.
//
func WithCancelChannel(parent context.Context, stopCh <-chan struct{}) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		select {
		case <-stopCh:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

// BindingAdapter receives from a binding.Receiver and sends to a binding.Sender.
type BindingAdapter struct {
	Receiver binding.Receiver
	Sender   binding.Sender
}

// Run runs a receive/send loop till there is an error, and returns the error.
//
// Returns nil if Receive() returns io.EOF, indicating an orderly shutdown.
func (p *BindingAdapter) Run(ctx context.Context) (err error) {
	defer func() {
		for _, i := range []interface{}{p.Receiver, p.Sender} {
			if c, _ := i.(binding.Closer); c != nil {
				if err2 := c.Close(context.Background()); err == nil {
					err = err2
				}
			}
		}
	}()

	for {
		m, err := p.Receiver.Receive(ctx)
		if err != nil {
			logging.FromContext(ctx).Error("receive error", zap.Error(err))
			return err
		}
		if err = p.Sender.Send(ctx, m); err != nil {
			logging.FromContext(ctx).Error("send error", zap.Error(err))
			return err
		}
	}
}

// Start implements adapter.Start
func (p *BindingAdapter) Start(stopCh <-chan struct{}) error {
	ctx, cancel := WithCancelChannel(context.Background(), stopCh)
	defer cancel() // Ensure no resource leaks.
	return p.Run(ctx)
}

// TODO(alanconway) following belong in cloudevents/sdk-go

// NewClientSender wraps a cloudevents.Client as a binding.Sender so
// we can use the default sink client provided by adapter.Main() as a Sender.
func NewClientSender(c client.Client, _ context.Context) binding.SendCloser {
	return &clientSender{client: c}
}

type clientSender struct{ client client.Client }

func (s *clientSender) Send(ctx context.Context, m binding.Message) (err error) {
	defer m.Finish(err)
	e, err := m.Event()
	if err != nil {
		return err
	}
	_, _, err = s.client.Send(ctx, e)
	return err
}

func (*clientSender) Close(_ context.Context) error { return nil }
