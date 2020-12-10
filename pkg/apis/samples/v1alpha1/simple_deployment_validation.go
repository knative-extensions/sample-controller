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

package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
)

// Validate implements apis.Validatable
func (d *SimpleDeployment) Validate(ctx context.Context) *apis.FieldError {
	return d.Spec.Validate(ctx).ViaField("spec")
}

// Validate implements apis.Validatable
func (ds *SimpleDeploymentSpec) Validate(ctx context.Context) *apis.FieldError {
	if ds.Image == "" {
		return apis.ErrMissingField("image")
	}
	return nil
}
