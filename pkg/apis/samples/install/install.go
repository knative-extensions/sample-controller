package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"knative.dev/sample-controller/pkg/apis/samples"
	"knative.dev/sample-controller/pkg/apis/samples/v1alpha1"
)

func Install(scheme *runtime.Scheme) {
	utilruntime.Must(samples.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	utilruntime.Must(scheme.SetVersionPriority(
		v1alpha1.SchemeGroupVersion,
	))
}
