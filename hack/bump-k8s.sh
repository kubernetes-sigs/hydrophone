#!/bin/bash

# Copyright 2024 The Kubernetes Authors.
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

set -xeuo pipefail

HYDROPHONE_ROOT=$(git rev-parse --show-toplevel)
echo "HYDROPHONE_ROOT: $HYDROPHONE_ROOT"

pushd "${HYDROPHONE_ROOT}" >/dev/null
  # Direct deps only: upgrading indirect k8s.io/* modules (kube-openapi, ...) past
  # what k8s itself pins can produce incompatible combinations that do not compile.
  go mod edit -json | jq -r ".Require[] | select(.Indirect | not) | .Path | select(contains(\"k8s.io/\"))" | xargs go get
  go mod tidy

  K8S_VERSION=$(curl -fsSL https://dl.k8s.io/release/stable.txt)
  sed -i -r "s/conformance:v[0-9]+\.[0-9]+\.[0-9]+/conformance:$K8S_VERSION/g" README.md docs/*.md

popd >/dev/null
git status
