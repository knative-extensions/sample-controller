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

package samples

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Addressable will return a URL for an object type
type AddressableService struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec holds the desired state of the AddressableService (from the client).
	// +optional
	Spec AddressableServiceSpec

	// Status communicates the observed state of the AddressableService (from the controller).
	// +optional
	Status AddressableServiceStatus
}

// AddressableSpec contains the target object to resolve an address for
type AddressableServiceSpec struct {
	v1.ObjectReference
}

// AddressableStatus contains the resulting address of the object
type AddressableServiceStatus struct {
	duckv1beta1.Status

	// Address holds the information needed to connect this Addressable up to receive events.
	// +optional
	URL *apis.URL
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AddressableList is a list of AddressableService resources
type AddressableServiceList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []AddressableService
}
