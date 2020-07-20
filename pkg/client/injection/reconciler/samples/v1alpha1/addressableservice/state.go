/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Veroute.on 2.0 (the "License");
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
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/reconciler"
	"knative.dev/sample-controller/pkg/apis/samples/v1alpha1"
)

// State is used to track the state of a reconciler in a single run.
type State struct {
	// Key is the original reconciliation key from the queue.
	Key string
	// Namespace is the namespace split from the reconciliation key.
	Namespace string
	// Namespace is the name split from the reconciliation key.
	Name string

	// reconciler is the reconciler.
	reconciler Interface

	// rof is the read only interface cast of the reconciler.
	roi ReadOnlyInterface
	// IsROI (Read Only Interface) the reconciler only observes reconciliation.
	isROI bool
	// rof is the read only finalizer cast of the reconciler.
	rof ReadOnlyFinalizer
	// IsROF (Read Only Finalizer) the reconciler only observes finalize.
	isROF bool
	// IsLeader the instance of the reconciler is the elected leader.
	isLeader bool
	// IsStatusUpdated the resource had its status updated successful.
	isStatusUpdated bool
}

func NewState(key string, reconciler Interface) (*State, error) {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return nil, fmt.Errorf("invalid resource key: %s", key)
	}

	roi, isROI := reconciler.(ReadOnlyInterface)
	rof, isROF := reconciler.(ReadOnlyFinalizer)

	return &State{
		Key:        key,
		Namespace:  namespace,
		Name:       name,
		reconciler: reconciler,
		roi:        roi,
		isROI:      isROI,
		rof:        rof,
		isROF:      isROF,
	}, nil
}

// NamespacedName is a helper to create a types.NamespacedName
func (s *State) NamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: s.Namespace,
		Name:      s.Name,
	}
}

// IsNOP checks to see if this reconciler with the current state is enabled to
// do any work or not.
// IsNOP returns true when there is no work possible for the reconciler.
func (s *State) IsNOP(isLeader bool) bool {
	s.isLeader = isLeader
	if !s.isLeader && !s.isROI && !s.isROF {
		// If we are not the leader, and we don't implement either ReadOnly
		// interface, then take a fast-path out.
		return true
	}
	return false
}

func (s *State) ShouldRecord(event *reconciler.ReconcilerEvent) bool {
	if event.ConditionFn == nil {
		return true
	}
	if event.ConditionFn(s) {
		return true
	}
	return false
}

// IsLeader Implements knative.dev/pkg/reconciler/State.IsLeader
func (s *State) IsLeader() bool {
	return s.isLeader
}

// StatusUpdated  marks the status as being updated for this state.
func (s *State) StatusUpdated() {
	s.isStatusUpdated = true
}

// IsStatusUpdated Implements knative.dev/pkg/reconciler/State.IsStatusUpdated
func (s *State) IsStatusUpdated() bool {
	return s.isStatusUpdated
}

func (s *State) ReconcileMethodFor(resource *v1alpha1.AddressableService) (string, DoReconcile) {
	if resource.GetDeletionTimestamp().IsZero() {
		if s.isLeader {
			return "ReconcileKind", s.reconciler.ReconcileKind
		} else if s.isROI {
			return "ObserveKind", s.roi.ObserveKind
		}
	} else if fin, ok := s.reconciler.(Finalizer); s.isLeader && ok {
		return "FinalizeKind", fin.FinalizeKind
	} else if !s.isLeader && s.isROF {
		return "ObserveFinalizeKind", s.rof.ObserveFinalizeKind
	}
	return "unknown", nil
}
