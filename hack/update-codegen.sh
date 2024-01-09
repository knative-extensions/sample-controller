#!/usr/bin/env bash

# Copyright 2019 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

source $(dirname $0)/../vendor/knative.dev/hack/codegen-library.sh
export PATH="$GOBIN:$PATH"

function run_yq() {
	go_run github.com/mikefarah/yq/v4@v4.23.1 "$@"
}

echo "=== Update Codegen for ${MODULE_NAME}"

group "Kubernetes Codegen"

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
${CODEGEN_PKG}/generate-groups.sh "deepcopy,client,informer,lister" \
  knative.dev/sample-controller/pkg/client knative.dev/sample-controller/pkg/apis \
  "samples:v1alpha1" \
  --go-header-file ${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt

group "Knative Codegen"

# Knative Injection
${KNATIVE_CODEGEN_PKG}/hack/generate-knative.sh "injection" \
  knative.dev/sample-controller/pkg/client knative.dev/sample-controller/pkg/apis \
  "samples:v1alpha1" \
  --go-header-file ${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt

group "Update CRD Schema"

go run $(dirname $0)/../cmd/schema/ dump SimpleDeployment \
  | run_yq eval-all --header-preprocess=false --inplace 'select(fileIndex == 0).spec.versions[0].schema.openAPIV3Schema = select(fileIndex == 1) | select(fileIndex == 0)' \
  $(dirname $0)/../config/300-simpledeployment.yaml -

go run $(dirname $0)/../cmd/schema/ dump AddressableService \
  | run_yq eval-all --header-preprocess=false --inplace 'select(fileIndex == 0).spec.versions[0].schema.openAPIV3Schema = select(fileIndex == 1) | select(fileIndex == 0)' \
  $(dirname $0)/../config/300-addressableservice.yaml -

group "Update deps post-codegen"

# Make sure our dependencies are up-to-date
${REPO_ROOT_DIR}/hack/update-deps.sh
