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

// Package adapter implements a sample receive adapter which receives events
// via HTTP POST requests.
package adapter

// This file implements the send/receive logic for the Adapter interface.

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/binding/format"
	"github.com/cloudevents/sdk-go/pkg/binding/spec"
	ce "github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"go.uber.org/zap"
	"knative.dev/eventing/pkg/adapter"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/source"
)

type envConfig struct {
	// EnvConfig is included in all receive adapters.
	// It includes a SinkURI field for outgoing HTTP events.
	adapter.EnvConfig

	// SourceURI provides the listening address for incoming HTTP events.
	SourceURI string `envconfig:"SOURCE_URI" required:"true"`
}

func NewEnv() adapter.EnvConfigAccessor { return &envConfig{} }

// Adapter receives events as a HTTP server, and sends them via
// a CloudEvents client.
type Adapter struct {
	listener net.Listener
	server   http.Server
	sink     client.Client
}

func (a *Adapter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	e, err := decode(req)
	if err == nil {
		_, _, err = a.sink.Send(req.Context(), e)
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// Start runs the adapter. It returns if stopCh is closed or the adapter closes
// due to some error.
func (a *Adapter) Start(stopCh <-chan struct{}) error {
	retCh := make(chan struct{}) // Closed when Serve() returns
	defer func() { close(retCh) }()
	go func() { // Watch stopCh and Shutdown the server if closed.
		select {
		case <-stopCh:
			a.server.Shutdown(context.Background())
		case <-retCh: // We are return for some other reason.
		}
	}()
	return a.server.Serve(a.listener)
}

func NewAdapter(ctx context.Context, aEnv adapter.EnvConfigAccessor, sink client.Client, reporter source.StatsReporter) adapter.Adapter {
	env := aEnv.(*envConfig) // Will always be our own envConfig type
	log := logging.FromContext(ctx)
	u, err := url.Parse(env.SourceURI)
	if err != nil {
		log.Fatal("invalid URI", zap.String("source", env.SourceURI))
	}
	a := &Adapter{sink: sink}
	a.server.Handler = a
	a.listener, err = net.Listen("tcp", u.Host)
	if err != nil {
		log.Fatal("listen error", zap.Error(err))
	}
	log.Info("Sample adapter listening",
		zap.String("source", env.SourceURI),
		zap.String("sink", env.SinkURI))
	return a
}

// CloudEvents spec simplifies support multiple spec. versions.
// HTTP uses the "ce-" prefix for CloudEvent attribute headers.
var specs = spec.WithPrefix("ce-")

// decode a http.Request as a CloudEvent.
func decode(req *http.Request) (ce.Event, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return ce.Event{}, err
	}
	if ct := req.Header.Get("Content-Type"); format.IsFormat(ct) {
		// Structured mode message, body is formatted event.
		var e ce.Event
		err = format.Unmarshal(ct, body, &e)
		if err != nil {
			return ce.Event{}, err
		}
		return e, nil
	}
	// Binary mode message, body is event data, headers contain event attributes.
	version, err := specs.FindVersion(func(k string) string {
		return req.Header.Get(strings.ToLower(k))
	})
	if err != nil {
		return ce.Event{}, err
	}
	c := version.NewContext()
	if err := c.SetDataContentType(req.Header.Get("Content-Type")); err != nil {
		return ce.Event{}, err
	}
	for k, v := range req.Header {
		if err := version.SetAttribute(c, k, v[0]); err != nil {
			return ce.Event{}, err
		}
	}
	if len(body) == 0 {
		return ce.Event{Data: nil, Context: c}, nil
	}
	return ce.Event{Data: body, DataEncoded: true, Context: c}, nil
}
