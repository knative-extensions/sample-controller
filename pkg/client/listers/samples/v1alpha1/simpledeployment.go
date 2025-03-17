/*
Copyright 2025 The Knative Authors

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
	samplesv1alpha1 "knative.dev/sample-controller/pkg/apis/samples/v1alpha1"
)

// SimpleDeploymentLister helps list SimpleDeployments.
// All objects returned here must be treated as read-only.
type SimpleDeploymentLister interface {
	// List lists all SimpleDeployments in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*samplesv1alpha1.SimpleDeployment, err error)
	// SimpleDeployments returns an object that can list and get SimpleDeployments.
	SimpleDeployments(namespace string) SimpleDeploymentNamespaceLister
	SimpleDeploymentListerExpansion
}

// simpleDeploymentLister implements the SimpleDeploymentLister interface.
type simpleDeploymentLister struct {
	listers.ResourceIndexer[*samplesv1alpha1.SimpleDeployment]
}

// NewSimpleDeploymentLister returns a new SimpleDeploymentLister.
func NewSimpleDeploymentLister(indexer cache.Indexer) SimpleDeploymentLister {
	return &simpleDeploymentLister{listers.New[*samplesv1alpha1.SimpleDeployment](indexer, samplesv1alpha1.Resource("simpledeployment"))}
}

// SimpleDeployments returns an object that can list and get SimpleDeployments.
func (s *simpleDeploymentLister) SimpleDeployments(namespace string) SimpleDeploymentNamespaceLister {
	return simpleDeploymentNamespaceLister{listers.NewNamespaced[*samplesv1alpha1.SimpleDeployment](s.ResourceIndexer, namespace)}
}

// SimpleDeploymentNamespaceLister helps list and get SimpleDeployments.
// All objects returned here must be treated as read-only.
type SimpleDeploymentNamespaceLister interface {
	// List lists all SimpleDeployments in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*samplesv1alpha1.SimpleDeployment, err error)
	// Get retrieves the SimpleDeployment from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*samplesv1alpha1.SimpleDeployment, error)
	SimpleDeploymentNamespaceListerExpansion
}

// simpleDeploymentNamespaceLister implements the SimpleDeploymentNamespaceLister
// interface.
type simpleDeploymentNamespaceLister struct {
	listers.ResourceIndexer[*samplesv1alpha1.SimpleDeployment]
}
