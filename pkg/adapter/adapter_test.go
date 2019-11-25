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
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	ce "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"knative.dev/eventing-contrib/pkg/kncloudevents"
	"knative.dev/eventing/pkg/adapter"
)

func TestAdapter(t *testing.T) {
	// Test sink to receive events.
	sink := newSink()
	defer sink.Close()

	// Run the adapter in this process.
	c, err := kncloudevents.NewDefaultClient(sink.URL)
	require.NoError(t, err)
	a := NewAdapter(context.Background(), &envConfig{SourceURI: "http://:0"}, c, nil)
	sourceURI := "http://" + a.(*Adapter).listener.Addr().String()
	stop := make(chan struct{})
	go a.Start(stop)
	defer func() { close(stop) }()

	c, err = kncloudevents.NewDefaultClient(sourceURI)
	require.NoError(t, err)
	_, _, err = c.Send(context.Background(), testEvent)
	assert.NoError(t, err)
	got := <-sink.received
	assert.Equal(t, testEvent, got)
}

func TestAdapterMain(t *testing.T) {
	// Use the test executable to simulate the cmd/receive_adapter process if
	// environment var t.Name() is set to "main"
	// (see https://talks.golang.org/2014/testing.slide#23)
	if os.Getenv(t.Name()) == "main" {
		adapter.Main("sample-source", NewEnv, NewAdapter)
		return
	}

	// Set up a test sink to receive from the adapter.
	received := make(chan ce.Event, 1)
	sink := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		e, err := decode(req)
		if err != nil {
			e.SetDataContentType(err.Error()) // invalidate the event with error message
		}
		received <- e
	}))
	defer sink.Close()
	sinkURI := sink.URL

	// Run a simulated receive_adapter main using the test executable.
	sourceURI := "http://" + ephemeralAddr(t)
	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	cmd.Env = append(os.Environ(),
		t.Name()+"=main",
		"SINK_URI="+sinkURI,
		"SOURCE_URI="+sourceURI,
		"NAMESPACE=namespace",
		`K_METRICS_CONFIG={"domain":"x", "component":"x", "prometheusport":0, "configmap":{}}`,
		`K_LOGGING_CONFIG={}`,
	)

	// Collect output, wait for "starting sample adapter"
	pr, pw := io.Pipe()
	cmd.Stdout, cmd.Stderr = pw, pw
	require.NoError(t, cmd.Start())
	defer func() { cmd.Process.Kill(); _, _ = cmd.Process.Wait() }()
	require.NoError(t, waitFor(pr, "Sample adapter listening"))
	go func() { io.Copy(cmd.Stderr, pr) }() // Copy any further log output to stderr.

	// Send to adapter source URI, verify event received at the sink.
	e := ce.Event{
		Context: ce.EventContextV1{
			Type:   "com.example.MinEvent",
			Source: types.URIRef{URL: url.URL{Scheme: "http", Host: "example.com", Path: "/source"}},
			ID:     "min-event",
			Time:   &types.Timestamp{Time: time.Date(2020, 03, 21, 12, 34, 56, 780000000, time.UTC)},
		}.AsV1(),
	}
	c, err := kncloudevents.NewDefaultClient(sourceURI)
	require.NoError(t, err)
	_, _, err = c.Send(context.Background(), e)
	assert.NoError(t, err)
	got := <-received
	assert.Equal(t, e, got)
}

var testEvent = ce.Event{
	Context: ce.EventContextV1{
		Type:   "com.example.MinEvent",
		Source: types.URIRef{URL: url.URL{Scheme: "http", Host: "example.com", Path: "/source"}},
		ID:     "min-event",
		Time:   &types.Timestamp{Time: time.Date(2020, 03, 21, 12, 34, 56, 780000000, time.UTC)},
	}.AsV1(),
}

func adapterEnv(sourceURI, sinkURI string) []string {
	return []string{
		"SINK_URI=" + sinkURI,
		"SOURCE_URI=" + sourceURI,
		"NAMESPACE=namespace",
		`K_METRICS_CONFIG={"domain":"x", "component":"x", "prometheusport":0, "configmap":{}}`,
		`K_LOGGING_CONFIG={}`,
	}
}

// ephemeralAddr returns a listening address with a free local ephemeral port.
//
// This is not 100% reliable, the port can be taken before the caller starts
// listening on it.
func ephemeralAddr(t *testing.T) string {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()
	return l.Addr().String()
}

func waitFor(r io.Reader, what string) error {
	out := ""
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		out += scanner.Text() + "\n"
		if strings.Contains(out, what) {
			return nil
		}
	}
	return fmt.Errorf("Expected: %q\nGot ==\n%s\n--", what, out)
}

type sink struct {
	*httptest.Server
	received chan ce.Event
}

func (s *sink) ServeHTTP(_ http.ResponseWriter, req *http.Request) {
	e, err := decode(req)
	if err != nil {
		e.SetDataContentType(err.Error()) // invalidate the event with error message
	}
	s.received <- e
}

func newSink() *sink {
	s := &sink{}
	s.received = make(chan ce.Event, 1) // Don't block
	s.Server = httptest.NewServer(s)
	return s
}
