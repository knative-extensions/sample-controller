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

// +genreconciler
// +genreconciler:kind=knative.dev/sample-controller/pkg/apis/samples/v1alpha1.AddressableService
// +genreconciler:injection-client=knative.dev/sample-controller/pkg/client/injection/client
// +genreconciler:injection-informer=knative.dev/sample-controller/pkg/client/injection/informers/samples/v1alpha1/addressableservice
// +genreconciler:clientset=knative.dev/sample-controller/pkg/client/clientset/versioned
// +genreconciler:lister=knative.dev/sample-controller/pkg/client/listers/samples/v1alpha1

package addressableservice
