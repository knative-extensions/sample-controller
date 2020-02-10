/*
Copyright 2020 The Knative Authors

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

// Code generated by injection-gen. DO NOT EDIT.

package addressableservice

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	watch "k8s.io/apimachinery/pkg/watch"
	scheme "k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	record "k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	client "knative.dev/pkg/client/injection/kube/client"
	controller "knative.dev/pkg/controller"
	logging "knative.dev/pkg/logging"
	versionedscheme "knative.dev/sample-controller/pkg/client/clientset/versioned/scheme"
	injectionclient "knative.dev/sample-controller/pkg/client/injection/client"
	addressableservice "knative.dev/sample-controller/pkg/client/injection/informers/samples/v1alpha1/addressableservice"
)

const (
	defaultControllerAgentName = "addressableservice-controller"
	defaultFinalizerName       = "addressableservice" // TODO: make this have the api group.
)

func NewFinalizingImpl(ctx context.Context, r Interface, finalizer string) *controller.Impl {
	logger := logging.FromContext(ctx)
	rec := newRecordedReconcilerImpl(ctx, r)
	finalizer = strings.TrimSpace(finalizer)
	if finalizer == "" {
		finalizer = defaultFinalizerName
	}
	rec.finalizerName = pointer.StringPtr(finalizer)
	return controller.NewImpl(rec, logger, "addressableservices")
}

func NewImpl(ctx context.Context, r Interface) *controller.Impl {
	logger := logging.FromContext(ctx)
	rec := newRecordedReconcilerImpl(ctx, r)
	return controller.NewImpl(rec, logger, "addressableservices")
}

func newRecordedReconcilerImpl(ctx context.Context, r Interface) *reconcilerImpl {
	logger := logging.FromContext(ctx)

	addressableserviceInformer := addressableservice.Get(ctx)

	recorder := controller.GetEventRecorder(ctx)
	if recorder == nil {
		// Create event broadcaster
		logger.Debug("Creating event broadcaster")
		eventBroadcaster := record.NewBroadcaster()
		watches := []watch.Interface{
			eventBroadcaster.StartLogging(logger.Named("event-broadcaster").Infof),
			eventBroadcaster.StartRecordingToSink(
				&v1.EventSinkImpl{Interface: client.Get(ctx).CoreV1().Events("")}),
		}
		recorder = eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: defaultControllerAgentName})
		go func() {
			<-ctx.Done()
			for _, w := range watches {
				w.Stop()
			}
		}()
	}

	return &reconcilerImpl{
		Client:     injectionclient.Get(ctx),
		Lister:     addressableserviceInformer.Lister(),
		Recorder:   recorder,
		reconciler: r,
	}
}

func init() {
	versionedscheme.AddToScheme(scheme.Scheme)
}
