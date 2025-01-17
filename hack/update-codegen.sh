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

echo "=== Update Codegen for ${MODULE_NAME}"

group "Kubernetes Codegen"

source "${CODEGEN_PKG}/kube_codegen.sh"

kube::codegen::gen_client \
  --boilerplate "${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt" \
  --output-dir "${REPO_ROOT_DIR}/pkg/client" \
  --output-pkg "knative.dev/sample-controller/pkg/client" \
  --with-watch \
  "${REPO_ROOT_DIR}/pkg/apis"

kube::codegen::gen_helpers \
  --boilerplate "${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt" \
  "${REPO_ROOT_DIR}/pkg"

group "Knative Codegen"

# Knative Injection
${KNATIVE_CODEGEN_PKG}/hack/generate-knative.sh "injection" \
  knative.dev/sample-controller/pkg/client knative.dev/sample-controller/pkg/apis \
  "samples:v1alpha1" \
  --go-header-file ${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt

group "Update CRD Schema"

go run sigs.k8s.io/controller-tools/cmd/controller-gen@v0.17.1 \
  schemapatch:manifests=config,generateEmbeddedObjectMeta=true \
  output:dir=config/ \
  paths=./pkg/apis/...

group "Update deps post-codegen"
# Make sure our dependencies are up-to-date
${REPO_ROOT_DIR}/hack/update-deps.sh
