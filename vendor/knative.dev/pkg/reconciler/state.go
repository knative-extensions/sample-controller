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

package reconciler

import (
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
)

// State is used to track the state of a reconciler in a single run.
type State struct {
	// Key is the original reconciliation key from the queue.
	Key string
	// Namespace is the namespace split from the reconciliation key.
	Namespace string
	// Namespace is the name split from the reconciliation key.
	Name string

	// IsROI (Read Only Interface) the reconciler only observes reconciliation.
	IsROI bool
	// IsROF (Read Only Finalizer) the reconciler only observes finalize.
	IsROF bool
	// IsLeader the instance of the reconciler is the elected leader.
	IsLeader bool
	// IsStatusUpdated the resource had its status updated successful.
	IsStatusUpdated bool
}

func NewState(key string) (*State, error) {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return nil, fmt.Errorf("invalid resource key: %s", key)
	}
	return &State{
		Key:       key,
		Namespace: namespace,
		Name:      name,
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
func (s *State) IsNOP() bool {
	if !s.IsLeader && !s.IsROI && !s.IsROF {
		return true
	}
	return false
}

func (s *State) ShouldRecord(event *ReconcilerEvent) bool {
	if event.ConditionFn == nil {
		return true
	}
	if event.ConditionFn(s) {
		return true
	}
	return false
}
