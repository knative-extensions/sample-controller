package v1beta1

import (
	conversion "k8s.io/apimachinery/pkg/conversion"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
	"knative.dev/sample-controller/pkg/apis/samples"
)

func Convert_v1beta1_AddressableServiceStatus_To_samples_AddressableServiceStatus(
	in *AddressableServiceStatus,
	out *samples.AddressableServiceStatus,
	s conversion.Scope,
) error {

	if in.Address != nil {
		out.URL = in.Address.URL
	}
	return nil
}

func Convert_samples_AddressableServiceStatus_To_v1beta1_AddressableServiceStatus(
	in *samples.AddressableServiceStatus,
	out *AddressableServiceStatus,
	s conversion.Scope,
) error {

	out.Address = &duckv1beta1.Addressable{
		URL: in.URL,
	}

	return nil
}
