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
	"fmt"
	"io"
	"io/ioutil"
	nethttp "net/http"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/binding/format"
	"github.com/cloudevents/sdk-go/pkg/binding/spec"
	ce "github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
)

// Prefix defined by the CloudEvents HTTP binding.
const prefix = "ce-"

// Helper to decode CloudEvents encoded using any knonw spec-version.
var specs = spec.WithPrefix(prefix)

const ContentType = "Content-Type"

// Message holds the Header and Body of a HTTP Request or Response.
// and implements the binding.Message interface to decode and
// forward events.
type Message struct {
	Header   nethttp.Header
	Body     []byte
	OnFinish func(error) error
}

func NewMessage(header nethttp.Header, body io.ReadCloser) (*Message, error) {
	m := Message{Header: header}
	if body != nil {
		defer func() { _ = body.Close() }()
		var err error
		if m.Body, err = ioutil.ReadAll(body); err != nil && err != io.EOF {
			return nil, err
		}
		if len(m.Body) == 0 {
			m.Body = nil
		}
	}
	return &m, nil
}

// Structured returns structured content if the HTTP request was encoded in structured mode.
func (m *Message) Structured() (string, []byte) {
	if ct := m.Header.Get(ContentType); format.IsFormat(ct) {
		return ct, m.Body
	}
	return "", nil
}

// Find the CloudEvents spec version for a HTTP message.
func (m *Message) findVersion() (spec.Version, error) {
	for _, sv := range specs.SpecVersionNames() {
		if v, err := specs.Version(m.Header.Get(sv)); err == nil {
			return v, nil
		}
	}
	return nil, fmt.Errorf("CloudEvents spec-version not found")
}

// Set an EventContextWriter attribute from a HTTP header.
func (m *Message) setAttribute(version spec.Version, c ce.EventContextWriter, k string, v interface{}) error {
	k = strings.ToLower(k)
	if a := version.Attribute(k); a != nil { // Standard attribute
		return a.Set(c, v)
	}
	var err error
	if strings.HasPrefix(k, prefix) { // Extension attribute
		v, err = types.Validate(v)
		if err == nil {
			err = c.SetExtension(strings.TrimPrefix(k, prefix), v)
		}
	}
	return err
}

// Event decodes the contained event.
func (m *Message) Event() (e ce.Event, err error) {
	if f, b := m.Structured(); f != "" {
		err := format.Unmarshal(f, b, &e)
		return e, err
	}
	version, err := m.findVersion()
	if err != nil {
		return e, err
	}
	c := version.NewContext()
	if err := c.SetDataContentType(m.Header.Get(ContentType)); err != nil {
		return e, err
	}
	for k, v := range m.Header {
		if err := m.setAttribute(version, c, k, v[0]); err != nil {
			return e, err
		}
	}
	if len(m.Body) == 0 {
		return ce.Event{Data: nil, Context: c}, nil
	}
	return ce.Event{Data: m.Body, DataEncoded: true, Context: c}, nil
}

// Finish indicates the sender is finished with the message.
func (m *Message) Finish(err error) error {
	if m.OnFinish != nil {
		return m.OnFinish(err)
	}
	return nil
}
