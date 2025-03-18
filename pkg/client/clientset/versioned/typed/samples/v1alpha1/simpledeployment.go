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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
	samplesv1alpha1 "knative.dev/sample-controller/pkg/apis/samples/v1alpha1"
	scheme "knative.dev/sample-controller/pkg/client/clientset/versioned/scheme"
)

// SimpleDeploymentsGetter has a method to return a SimpleDeploymentInterface.
// A group's client should implement this interface.
type SimpleDeploymentsGetter interface {
	SimpleDeployments(namespace string) SimpleDeploymentInterface
}

// SimpleDeploymentInterface has methods to work with SimpleDeployment resources.
type SimpleDeploymentInterface interface {
	Create(ctx context.Context, simpleDeployment *samplesv1alpha1.SimpleDeployment, opts v1.CreateOptions) (*samplesv1alpha1.SimpleDeployment, error)
	Update(ctx context.Context, simpleDeployment *samplesv1alpha1.SimpleDeployment, opts v1.UpdateOptions) (*samplesv1alpha1.SimpleDeployment, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, simpleDeployment *samplesv1alpha1.SimpleDeployment, opts v1.UpdateOptions) (*samplesv1alpha1.SimpleDeployment, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*samplesv1alpha1.SimpleDeployment, error)
	List(ctx context.Context, opts v1.ListOptions) (*samplesv1alpha1.SimpleDeploymentList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *samplesv1alpha1.SimpleDeployment, err error)
	SimpleDeploymentExpansion
}

// simpleDeployments implements SimpleDeploymentInterface
type simpleDeployments struct {
	*gentype.ClientWithList[*samplesv1alpha1.SimpleDeployment, *samplesv1alpha1.SimpleDeploymentList]
}

// newSimpleDeployments returns a SimpleDeployments
func newSimpleDeployments(c *SamplesV1alpha1Client, namespace string) *simpleDeployments {
	return &simpleDeployments{
		gentype.NewClientWithList[*samplesv1alpha1.SimpleDeployment, *samplesv1alpha1.SimpleDeploymentList](
			"simpledeployments",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *samplesv1alpha1.SimpleDeployment { return &samplesv1alpha1.SimpleDeployment{} },
			func() *samplesv1alpha1.SimpleDeploymentList { return &samplesv1alpha1.SimpleDeploymentList{} },
		),
	}
}
