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

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/binding"
	"github.com/cloudevents/sdk-go/pkg/binding/test"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/stretchr/testify/require"
	"knative.dev/eventing-contrib/pkg/kncloudevents"
	"knative.dev/sample-controller/pkg/eventing/adapter"
)

var (
	sourceURI = types.URIRef{URL: url.URL{Scheme: "http", Host: "example.com", Path: "/source"}}
	timestamp = types.Timestamp{Time: time.Date(2020, 03, 21, 12, 34, 56, 780000000, time.UTC)}
	schema    = types.URI{URL: url.URL{Scheme: "http", Host: "example.com", Path: "/schema"}}
)

func testEvent() cloudevents.Event {
	return cloudevents.Event{
		Context: cloudevents.EventContextV1{
			Type:   "com.example.MinEvent",
			Source: sourceURI,
			ID:     "min-event",
			Time:   &timestamp,
		}.AsV1(),
	}
}

func testSender(t *testing.T, targetURL string) binding.SendCloser {
	t.Helper()
	c, err := kncloudevents.NewDefaultClient(targetURL)
	require.NoError(t, err)
	return adapter.NewClientSender(c, context.Background())
}

func TestNewHandlerReceiver(t *testing.T) {
	r := NewHandlerReceiver()
	rSrv := httptest.NewServer(r)
	defer rSrv.Close()
	s := testSender(t, rSrv.URL)
	defer s.Close(context.Background())

	want := binding.EventMessage(testEvent())
	got := test.SendReceive(t, want, s, r)
	test.AssertMessageEventEqual(t, want, got)
}

func TestNewServerReceiver(t *testing.T) {
	addr := ephemeralAddr(t)
	r, err := NewServerReceiver(http.Server{Addr: addr})
	require.NoError(t, err)
	defer r.Close(context.Background())
	s := testSender(t, "http://"+addr)
	defer s.Close(context.Background())

	want := binding.EventMessage(testEvent())
	got := test.SendReceive(t, want, s, r)
	test.AssertMessageEventEqual(t, want, got)
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

func TestAdapterMain(t *testing.T) {
	// Use the test executable to simulate the cmd/receive_adapter process if env var
	// t.Name() is set to "main" (trick from https://talks.golang.org/2014/testing.slide#23)
	if os.Getenv(t.Name()) == "main" {
		adapter.Main("sample-source", NewEnv, NewAdapter)
		return
	}

	// Use an adapter.Receiver as the sink for the adapter.
	r := NewHandlerReceiver()
	rSrv := httptest.NewServer(r)
	defer rSrv.Close()
	sinkURI := rSrv.URL

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
	require.NoError(t, waitFor(pr, "starting sample receive adapter"))
	go func() { io.Copy(cmd.Stderr, pr) }() // Copy any further log output to stderr.

	// Send to adapter source URI, verify event received at the sink.
	s := testSender(t, sourceURI)
	defer s.Close(context.Background())
	want := binding.EventMessage(testEvent())
	got := test.SendReceive(t, want, s, r)
	test.AssertMessageEventEqual(t, want, got)
}
