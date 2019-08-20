package samples

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
)

var condSet = apis.NewLivingConditionSet()

const (
	// AddressableServiceConditionReady is set when the revision is starting to materialize
	// runtime resources, and becomes true when those resources are ready.
	AddressableServiceConditionReady = apis.ConditionReady
)

// GetGroupVersionKind implements kmeta.OwnerRefable
func (as *AddressableService) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("AddressableService")
}

func (ass *AddressableServiceStatus) InitializeConditions() {
	condSet.Manage(ass).InitializeConditions()
}

func (ass *AddressableServiceStatus) MarkServiceUnavailable(name string) {
	condSet.Manage(ass).MarkFalse(
		AddressableServiceConditionReady,
		"ServiceUnavailable",
		"Service %q wasn't found.", name)
}

func (ass *AddressableServiceStatus) MarkServiceAvailable() {
	condSet.Manage(ass).MarkTrue(AddressableServiceConditionReady)
}
