/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"context"

	v1beta1 "k8s.io/api/admissionregistration/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	admissionregistrationv1beta1 "k8s.io/client-go/applyconfigurations/admissionregistration/v1beta1"
	gentype "k8s.io/client-go/gentype"
	scheme "k8s.io/client-go/kubernetes/scheme"
)

// ValidatingWebhookConfigurationsGetter has a method to return a ValidatingWebhookConfigurationInterface.
// A group's client should implement this interface.
type ValidatingWebhookConfigurationsGetter interface {
	ValidatingWebhookConfigurations() ValidatingWebhookConfigurationInterface
}

// ValidatingWebhookConfigurationInterface has methods to work with ValidatingWebhookConfiguration resources.
type ValidatingWebhookConfigurationInterface interface {
	Create(ctx context.Context, validatingWebhookConfiguration *v1beta1.ValidatingWebhookConfiguration, opts v1.CreateOptions) (*v1beta1.ValidatingWebhookConfiguration, error)
	Update(ctx context.Context, validatingWebhookConfiguration *v1beta1.ValidatingWebhookConfiguration, opts v1.UpdateOptions) (*v1beta1.ValidatingWebhookConfiguration, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.ValidatingWebhookConfiguration, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.ValidatingWebhookConfigurationList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.ValidatingWebhookConfiguration, err error)
	Apply(ctx context.Context, validatingWebhookConfiguration *admissionregistrationv1beta1.ValidatingWebhookConfigurationApplyConfiguration, opts v1.ApplyOptions) (result *v1beta1.ValidatingWebhookConfiguration, err error)
	ValidatingWebhookConfigurationExpansion
}

// validatingWebhookConfigurations implements ValidatingWebhookConfigurationInterface
type validatingWebhookConfigurations struct {
	*gentype.ClientWithListAndApply[*v1beta1.ValidatingWebhookConfiguration, *v1beta1.ValidatingWebhookConfigurationList, *admissionregistrationv1beta1.ValidatingWebhookConfigurationApplyConfiguration]
}

// newValidatingWebhookConfigurations returns a ValidatingWebhookConfigurations
func newValidatingWebhookConfigurations(c *AdmissionregistrationV1beta1Client) *validatingWebhookConfigurations {
	return &validatingWebhookConfigurations{
		gentype.NewClientWithListAndApply[*v1beta1.ValidatingWebhookConfiguration, *v1beta1.ValidatingWebhookConfigurationList, *admissionregistrationv1beta1.ValidatingWebhookConfigurationApplyConfiguration](
			"validatingwebhookconfigurations",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *v1beta1.ValidatingWebhookConfiguration { return &v1beta1.ValidatingWebhookConfiguration{} },
			func() *v1beta1.ValidatingWebhookConfigurationList {
				return &v1beta1.ValidatingWebhookConfigurationList{}
			}),
	}
}
