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
	corev1listers "k8s.io/client-go/listers/core/v1"

	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/network"
	"knative.dev/pkg/reconciler"
	"knative.dev/pkg/tracker"
	samplesv1alpha1 "knative.dev/sample-controller/pkg/apis/samples/v1alpha1"
	addressableservicereconciler "knative.dev/sample-controller/pkg/client/injection/reconciler/samples/v1alpha1/addressableservice"
)

// Reconciler implements addressableservicereconciler.Interface for
// AddressableService resources.
type Reconciler struct {
	// Tracker builds an index of what resources are watching other resources
	// so that we can immediately react to changes tracked resources.
	Tracker tracker.Interface

	// Listers index properties about resources
	ServiceLister corev1listers.ServiceLister
}

// Check that our Reconciler implements Interface
var _ addressableservicereconciler.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, o *samplesv1alpha1.AddressableService) reconciler.Event {
	logger := logging.FromContext(ctx)

	if err := r.Tracker.TrackReference(tracker.Reference{
		APIVersion: "v1",
		Kind:       "Service",
		Name:       o.Spec.ServiceName,
		Namespace:  o.Namespace,
	}, o); err != nil {
		logger.Errorf("Error tracking service %s: %v", o.Spec.ServiceName, err)
		return err
	}

	if _, err := r.ServiceLister.Services(o.Namespace).Get(o.Spec.ServiceName); apierrs.IsNotFound(err) {
		logger.Info("Service does not yet exist:", o.Spec.ServiceName)
		o.Status.MarkServiceUnavailable(o.Spec.ServiceName)
		return nil
	} else if err != nil {
		logger.Errorf("Error reconciling service %s: %v", o.Spec.ServiceName, err)
		return err
	}

	o.Status.MarkServiceAvailable()
	o.Status.Address = &duckv1.Addressable{
		URL: &apis.URL{
			Scheme: "http",
			Host:   network.GetServiceHostname(o.Spec.ServiceName, o.Namespace),
		},
	}

	return nil
}
