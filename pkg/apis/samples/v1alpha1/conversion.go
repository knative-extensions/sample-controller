package v1alpha1

import (
	conversion "k8s.io/apimachinery/pkg/conversion"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
	"knative.dev/sample-controller/pkg/apis/samples"
)

func Convert_v1alpha1_AddressableServiceSpec_To_samples_AddressableServiceSpec(
	in *AddressableServiceSpec,
	out *samples.AddressableServiceSpec,
	s conversion.Scope,
) error {

	out.APIVersion = "v1"
	out.Kind = "Service"
	out.Name = in.ServiceName

	return nil
}

func Convert_samples_AddressableServiceSpec_To_v1alpha1_AddressableServiceSpec(
	in *samples.AddressableServiceSpec,
	out *AddressableServiceSpec,
	s conversion.Scope,
) error {

	out.ServiceName = in.Name

	return nil
}

func Convert_v1alpha1_AddressableServiceStatus_To_samples_AddressableServiceStatus(
	in *AddressableServiceStatus,
	out *samples.AddressableServiceStatus,
	s conversion.Scope,
) error {

	if in.Address != nil {
		out.URL = in.Address.URL
	}
	return nil
}

func Convert_samples_AddressableServiceStatus_To_v1alpha1_AddressableServiceStatus(
	in *samples.AddressableServiceStatus,
	out *AddressableServiceStatus,
	s conversion.Scope,
) error {

	out.Address = &duckv1beta1.Addressable{
		URL: in.URL,
	}

	return nil
}
