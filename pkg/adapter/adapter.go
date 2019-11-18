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

// Package adapter implements a sample receive adapter to receive HTTP events.
package adapter

import (
	"context"
	"net/http"
	"net/url"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"go.uber.org/zap"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/source"
	"knative.dev/sample-controller/pkg/eventing/adapter"
)

type envConfig struct {
	// Include the standard adapter.EnvConfig used by all adapters.
	adapter.EnvConfig

	// SourceURI provides the "host:port" listening address for the HTTP Source.
	SourceURI string `envconfig:"SOURCE_URI" required:"true"`
}

func NewEnv() adapter.EnvConfigAccessor { return &envConfig{} }

// NOTE: The sink argument is provided by the current adapter.Main.  If we want
// to allow arbitrary sinks, we should instead allow the adapter constructor to
// create its own sink Sender. The library can provide `DefaultHTTPSink()
// binding.Sender` for sources that want the default HTTP behavior.

func NewAdapter(ctx context.Context, envAcc adapter.EnvConfigAccessor, sink client.Client, reporter source.StatsReporter) adapter.Adapter {
	env := envAcc.(*envConfig)
	log := logging.FromContext(ctx)
	u, err := url.Parse(env.SourceURI)
	if err != nil {
		log.Fatal("invalid URI", zap.String("source", env.SourceURI))
	}
	log.Info("starting sample adapter",
		zap.String("source", env.SourceURI),
		zap.String("sink", env.SinkURI))
	sender := adapter.NewClientSender(sink, ctx)
	receiver, err := NewServerReceiver(http.Server{Addr: u.Host})
	if err != nil {
		log.Fatal("can't create source receiver", zap.Error(err))
	}
	log.Info("starting sample receive adapter")
	return &adapter.BindingAdapter{Receiver: receiver, Sender: sender}
}
