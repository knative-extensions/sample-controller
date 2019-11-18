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

package adapter

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	nethttp "net/http"
	"sync"

	"github.com/cloudevents/sdk-go/pkg/binding"
	"go.uber.org/zap"
	"knative.dev/pkg/logging"
)

// Receive HTTP requests as CloudEvents
type Receiver struct {
	incoming  chan *Message
	closeOnce sync.Once
	err       error
	server    *http.Server
	busy      sync.WaitGroup
}

// Server is nil unless r was created with NewServerReceiver
func (r *Receiver) Server() *http.Server { return r.server }

// ServeHTTP implements http.Handler. Blocks until Message.Finish() is called.
func (r *Receiver) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log := logging.FromContext(req.Context())
	var err error
	defer func() {
		if err != nil {
			msg := "cannot forward CloudEvent"
			log.Error(msg, zap.Error(err))
			nethttp.Error(rw, fmt.Sprintf("%s: %s", msg, err), http.StatusInternalServerError)
			r.closeErr(req.Context(), err)
		}
	}()
	var m *Message
	m, err = NewMessage(req.Header, req.Body)
	if err != nil {
		return
	}
	done := make(chan error)
	m.OnFinish = func(err error) error { done <- err; return nil }
	r.incoming <- m // Send to Receive
	err = <-done    // Wait for Message.Finish()
}

// NewHandlerReceiver creates a receiver that can be used as a http.Handler
// with a http.Server. It does not start an HTTP server, for that see NewServerReceiver.
func NewHandlerReceiver() *Receiver { return &Receiver{incoming: make(chan *Message)} }

// NewServerReceiver creates a receiver with it's own concurrent HTTP server.
//
// On return there is a server listening on serverConfig.Addr, it will be closed
// when the receiver closes.
//
func NewServerReceiver(serverConfig http.Server) (*Receiver, error) {
	r := NewHandlerReceiver()
	r.server = &serverConfig
	r.server.Handler = r
	// We don't actually use ListenAndServe as we want to be listening before we return.
	if r.server.Addr == "" {
		r.server.Addr = ":http"
	}
	l, err := net.Listen("tcp", r.server.Addr)
	if err != nil {
		return nil, err
	}
	r.busy.Add(1)
	go func() {
		defer r.busy.Done()
		_ = r.server.Serve(l)
	}()
	return r, nil
}

// Receive the next incoming HTTP request as a CloudEvent.
func (r *Receiver) Receive(ctx context.Context) (binding.Message, error) {
	if r.err != nil {
		return nil, r.err
	}
	m, ok := <-r.incoming
	if !ok {
		return nil, r.err
	}
	return m, nil
}

func (r *Receiver) closeErr(ctx context.Context, err error) error {
	r.closeOnce.Do(func() {
		if r.server != nil {
			if err2 := r.server.Shutdown(ctx); err == nil {
				err = err2
			}
			r.busy.Wait()
		}
		r.err = err
		if r.err == nil {
			r.err = io.EOF // r.err must be non-nil to signal closed
		}
		close(r.incoming)
	})
	if r.err == io.EOF { // closed cleanly, not an error for Close()
		return nil
	}
	return r.err
}

func (r *Receiver) Close(ctx context.Context) error { return r.closeErr(ctx, nil) }
