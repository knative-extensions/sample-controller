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

# shellcheck disable=SC1090
source "$(GOFLAGS='-mod=mod' go run knative.dev/hack/cmd/script codegen-library.sh)"

USAGE=$(cat <<EOF
Add boilerplate.<ext>.txt to all .<ext> files missing it in a directory.

Usage: (from repository root)
       ./hack/boilerplate/add-boilerplate.sh <ext> <DIR>

Example: (from repository root)
         ./hack/boilerplate/add-boilerplate.sh go cmd

As of now, only .go files are supported.
EOF
)

if [ -z "${1:-}" ] || [ -z "${2:-}" ]; then
  error Invalid arguments
  echo "${USAGE}" 1>&2
  exit 1
fi

if ! [[ "$1" = "go" ]]; then
  error Unsupported file extension
  echo "${USAGE}" 1>&2
  exit 2
fi

cnt=0
while read -r file; do
  if grep -q -E "^Copyright [[:digit:]]+ The Knative Authors$" "$file"; then
    continue
  fi
  cat "$(boilerplate)" > "$file".bck
  echo '' >> "$file".bck
  cat "$file" >> "$file".bck
  mv "$file".bck "$file"
  log License added to "$file"
  cnt=$(( cnt + 1))
done < <(find "$2" -type f -name "*.$1")

log License added to $cnt files
