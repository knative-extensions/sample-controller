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

# This script runs the presubmit tests; it is started by prow for each PR.
# For convenience, it can also be executed manually.
# Running the script without parameters, or with the --all-tests
# flag, causes all tests to be executed, in the right order.
# Use the flags --build-tests, --unit-tests and --integration-tests
# to run a specific set of tests.

# Markdown linting failures don't show up properly in Gubernator resulting
# in a net-negative contributor experience.
export DISABLE_MD_LINTING=1

source vendor/knative.dev/test-infra/scripts/presubmit-tests.sh


#/ We use the default build, unit test runners.
readonly INSTALL_YAML=$(mktemp)

function wait_for_result() {
  echo
  echo -n "Waiting for $1 reconciliation to run to completion."
  echo
  for i in {1..150}; do
    local url=$(kubectl get asvc my-address -o jsonpath="{.status.address.url}")

    if [[ "$url" == "http://my-service.default.svc.cluster.local" ]]; then
      echo
      echo "$1 tests passed"
      echo
      return 0
    fi

    echo -n "."
    sleep 2
  done

  echo "timed out waiting for result"
  return 1
}

function pre_integration_tests() {
  ko resolve -f config > "$INSTALL_YAML" || return 1
}

function post_integration_tests() {
  # Delete control plane
  kubectl delete -f "$INSTALL_YAML" || return 1

  # Delete CRDS
  kubectl delete \
    -f "config/v1alpha1" \
    -f "config/v1beta1" || return 1
}

function test_version() {
  version="$1"
  shift

  post_integration_tests

  # Install CRDs
  kubectl apply -f "config/$version" || return 1

  # Install control plane
  kubectl apply -f "$INSTALL_YAML" || return 1

  sleep 5s

  kubectl apply -f "test/$version.yaml" || return 1

  wait_for_result $version

  result=$?

  kubectl delete -f "test/$version.yaml"

  return $result
}

function integration_tests() {
  test_version v1alpha1 \
    || return 1

  test_version v1beta1 \
    || return 1
}


main $@
