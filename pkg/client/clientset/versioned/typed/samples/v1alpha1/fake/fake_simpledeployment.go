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

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "knative.dev/sample-controller/pkg/apis/samples/v1alpha1"
)

// FakeSimpleDeployments implements SimpleDeploymentInterface
type FakeSimpleDeployments struct {
	Fake *FakeSamplesV1alpha1
	ns   string
}

var simpledeploymentsResource = v1alpha1.SchemeGroupVersion.WithResource("simpledeployments")

var simpledeploymentsKind = v1alpha1.SchemeGroupVersion.WithKind("SimpleDeployment")

// Get takes name of the simpleDeployment, and returns the corresponding simpleDeployment object, and an error if there is any.
func (c *FakeSimpleDeployments) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.SimpleDeployment, err error) {
	emptyResult := &v1alpha1.SimpleDeployment{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(simpledeploymentsResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.SimpleDeployment), err
}

// List takes label and field selectors, and returns the list of SimpleDeployments that match those selectors.
func (c *FakeSimpleDeployments) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.SimpleDeploymentList, err error) {
	emptyResult := &v1alpha1.SimpleDeploymentList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(simpledeploymentsResource, simpledeploymentsKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SimpleDeploymentList{ListMeta: obj.(*v1alpha1.SimpleDeploymentList).ListMeta}
	for _, item := range obj.(*v1alpha1.SimpleDeploymentList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested simpleDeployments.
func (c *FakeSimpleDeployments) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(simpledeploymentsResource, c.ns, opts))

}

// Create takes the representation of a simpleDeployment and creates it.  Returns the server's representation of the simpleDeployment, and an error, if there is any.
func (c *FakeSimpleDeployments) Create(ctx context.Context, simpleDeployment *v1alpha1.SimpleDeployment, opts v1.CreateOptions) (result *v1alpha1.SimpleDeployment, err error) {
	emptyResult := &v1alpha1.SimpleDeployment{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(simpledeploymentsResource, c.ns, simpleDeployment, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.SimpleDeployment), err
}

// Update takes the representation of a simpleDeployment and updates it. Returns the server's representation of the simpleDeployment, and an error, if there is any.
func (c *FakeSimpleDeployments) Update(ctx context.Context, simpleDeployment *v1alpha1.SimpleDeployment, opts v1.UpdateOptions) (result *v1alpha1.SimpleDeployment, err error) {
	emptyResult := &v1alpha1.SimpleDeployment{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(simpledeploymentsResource, c.ns, simpleDeployment, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.SimpleDeployment), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSimpleDeployments) UpdateStatus(ctx context.Context, simpleDeployment *v1alpha1.SimpleDeployment, opts v1.UpdateOptions) (result *v1alpha1.SimpleDeployment, err error) {
	emptyResult := &v1alpha1.SimpleDeployment{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(simpledeploymentsResource, "status", c.ns, simpleDeployment, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.SimpleDeployment), err
}

// Delete takes name of the simpleDeployment and deletes it. Returns an error if one occurs.
func (c *FakeSimpleDeployments) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(simpledeploymentsResource, c.ns, name, opts), &v1alpha1.SimpleDeployment{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSimpleDeployments) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(simpledeploymentsResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.SimpleDeploymentList{})
	return err
}

// Patch applies the patch and returns the patched simpleDeployment.
func (c *FakeSimpleDeployments) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.SimpleDeployment, err error) {
	emptyResult := &v1alpha1.SimpleDeployment{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(simpledeploymentsResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.SimpleDeployment), err
}
