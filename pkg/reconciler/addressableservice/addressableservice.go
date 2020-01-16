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

package addressableservice

import (
	"context"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/network"
	"knative.dev/pkg/reconciler"
	"knative.dev/pkg/tracker"
	"knative.dev/sample-controller/pkg/apis/samples/v1alpha1"
)

// ReconcileKind implements Interface
func (r *Reconciler) ReconcileKind(ctx context.Context, asvc *v1alpha1.AddressableService) reconciler.Event {
	if asvc.GetDeletionTimestamp() != nil {
		// Check for a DeletionTimestamp.  If present, elide the normal reconcile logic.
		// When a controller needs finalizer handling, it would go here.
		return nil
	}
	asvc.Status.InitializeConditions()

	if err := r.reconcileService(ctx, asvc); err != nil {
		return err
	}

	asvc.Status.ObservedGeneration = asvc.Generation
	return nil
}

func (r *Reconciler) reconcileService(ctx context.Context, asvc *v1alpha1.AddressableService) error {
	logger := logging.FromContext(ctx)

	if err := r.Tracker.TrackReference(tracker.Reference{
		APIVersion: "v1",
		Kind:       "Service",
		Name:       asvc.Spec.ServiceName,
		Namespace:  asvc.Namespace,
	}, asvc); err != nil {
		return NewWarnInternal("Error tracking service %s: %v", asvc.Spec.ServiceName, err)
	}

	_, err := r.ServiceLister.Services(asvc.Namespace).Get(asvc.Spec.ServiceName)
	if apierrs.IsNotFound(err) {
		logger.Info("Service does not yet exist:", asvc.Spec.ServiceName)
		asvc.Status.MarkServiceUnavailable(asvc.Spec.ServiceName)
		return nil
	} else if err != nil {
		return NewWarnInternal("Error reconciling service %s: %v", asvc.Spec.ServiceName, err)
	}

	asvc.Status.MarkServiceAvailable()
	asvc.Status.Address = &duckv1.Addressable{
		URL: &apis.URL{
			Scheme: "http",
			Host:   network.GetServiceHostname(asvc.Spec.ServiceName, asvc.Namespace),
		},
	}
	return nil
}
