// +build !ignore_autogenerated

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

// Code generated by ___go_build_main_go. DO NOT EDIT.

package addressableservice

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	controller "knative.dev/pkg/controller"
	logging "knative.dev/pkg/logging"
	"knative.dev/pkg/tracker"
	client "knative.dev/sample-controller/pkg/client/injection/client"
	addressableservice "knative.dev/sample-controller/pkg/client/injection/informers/samples/v1alpha1/addressableservice"
)

const (
	controllerAgentName = "addressableservice-controller"
	finalizerName       = "addressableservice"
)

func NewImpl(ctx context.Context, r *Reconciler) *controller.Impl {
	logger := logging.FromContext(ctx)

	impl := controller.NewImpl(r, logger, "addressableservices")

	informer := addressableservice.Get(ctx)

	r.Core = Core{
		Client:  client.Get(ctx),
		Lister:  informer.Lister(),
		Tracker: tracker.New(impl.EnqueueKey, controller.GetTrackerLease(ctx)),
		Recorder: record.NewBroadcaster().NewRecorder(
			scheme.Scheme, v1.EventSource{Component: controllerAgentName}),
		FinalizerName: finalizerName,
		Reconciler:    r,
	}

	logger.Info("Setting up core event handlers")
	informer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	return impl
}
