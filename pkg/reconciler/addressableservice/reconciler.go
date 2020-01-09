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

package addressableservice

import (
	"context"
	"reflect"

	"knative.dev/pkg/reconciler"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/tracker"
	"knative.dev/sample-controller/pkg/apis/samples/v1alpha1"
	clientset "knative.dev/sample-controller/pkg/client/clientset/versioned"
	listers "knative.dev/sample-controller/pkg/client/listers/samples/v1alpha1"
)

// Interface defines the strongly typed interfaces to be implemented by a
// controller reconciling v1alpha1.AddressableService.
type Interface interface {
	// ReconcileKind implements custom logic to reconcile v1alpha1.AddressableService. Any changes
	// to the objects .Status or .Finalizers will be propagated to the stored
	// object. It is recommended that implementors do not call any update calls
	// for the Kind inside of ReconcileKind, it is the responsibility of the core
	// controller to propagate those properties.
	ReconcileKind(context.Context, *v1alpha1.AddressableService) reconciler.Event
}

var _ Interface = (*Reconciler)(nil)

// NewWarnInternal makes a new reconciler event with event type Warning, and
// reason InternalError.
func NewWarnInternal(msgf string, args ...interface{}) reconciler.Event {
	return reconciler.NewEvent(corev1.EventTypeWarning, "InternalError", msgf, args...)
}

// Reconciler implements controller.Reconciler for AddressableService resources.
type Reconciler struct {
	// Client is used to write back status updates.
	Client clientset.Interface

	// Listers index properties about resources
	Lister        listers.AddressableServiceLister
	ServiceLister corev1listers.ServiceLister

	// The tracker builds an index of what resources are watching other
	// resources so that we can immediately react to changes to changes in
	// tracked resources.
	Tracker tracker.Interface

	// Recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	Recorder record.EventRecorder
}

// Check that our Reconciler implements controller.Reconciler
var _ controller.Reconciler = (*Reconciler)(nil)

// Reconcile implements controller.Reconciler
func (r *Reconciler) Reconcile(ctx context.Context, key string) error {
	logger := logging.FromContext(ctx)

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		logger.Errorf("invalid resource key: %s", key)
		return nil
	}

	// If our controller has configuration state, we'd "freeze" it and
	// attach the frozen configuration to the context.
	//    ctx = r.configStore.ToContext(ctx)

	// Get the resource with this namespace/name.
	original, err := r.Lister.AddressableServices(namespace).Get(name)
	if apierrs.IsNotFound(err) {
		// The resource may no longer exist, in which case we stop processing.
		logger.Errorf("resource %q no longer exists", key)
		return nil
	} else if err != nil {
		return err
	}
	// Don't modify the informers copy.
	resource := original.DeepCopy()

	// Reconcile this copy of the resource and then write back any status
	// updates regardless of whether the reconciliation errored out.
	reconcileErr := r.ReconcileKind(ctx, resource)
	if equality.Semantic.DeepEqual(original.Status, resource.Status) {
		// If we didn't change anything then don't call updateStatus.
		// This is important because the copy we loaded from the informer's
		// cache may be stale and we don't want to overwrite a prior update
		// to status with this stale state.
	} else if _, err = r.updateStatus(resource); err != nil {
		logger.Warnw("Failed to update resource status", zap.Error(err))
		r.Recorder.Eventf(resource, corev1.EventTypeWarning, "UpdateFailed",
			"Failed to update status for %q: %v", resource.Name, err)
		return err
	}
	if reconcileErr != nil {
		logger.Error("ReconcileKind returned an error: %v", reconcileErr)
		var event *reconciler.ReconcilerEvent
		if reconciler.EventAs(reconcileErr, &event) {
			r.Recorder.Eventf(resource, event.EventType, event.Reason, event.Format, event.Args...)
		} else {
			r.Recorder.Event(resource, corev1.EventTypeWarning, "InternalError", reconcileErr.Error())
		}
	}
	return reconcileErr
}

// Update the Status of the resource.  Caller is responsible for checking
// for semantic differences before calling.
func (r *Reconciler) updateStatus(desired *v1alpha1.AddressableService) (*v1alpha1.AddressableService, error) {
	actual, err := r.Lister.AddressableServices(desired.Namespace).Get(desired.Name)
	if err != nil {
		return nil, err
	}
	// If there's nothing to update, just return.
	if reflect.DeepEqual(actual.Status, desired.Status) {
		return actual, nil
	}
	// Don't modify the informers copy
	existing := actual.DeepCopy()
	existing.Status = desired.Status
	return r.Client.SamplesV1alpha1().AddressableServices(desired.Namespace).UpdateStatus(existing)
}
