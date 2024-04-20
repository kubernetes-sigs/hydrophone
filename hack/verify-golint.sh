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

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE}")/..

cd "${KUBE_ROOT}"

LINT=${LINT:-golangci-lint}
VERSION=1.57.2

if [[ -z "$(command -v ${LINT})" ]]; then
  echo "${LINT} is missing. Installing it now..."

  LINT=$(go env GOPATH)/bin/golangci-lint
  mkdir -p "$(dirname "$LINT")"

  base="golangci-lint-$VERSION-$(go env GOOS)-$(go env GOARCH)"
  curl --fail -L https://github.com/golangci/golangci-lint/releases/download/v$VERSION/$base.tar.gz | tar xzOf - "$base/golangci-lint" > $LINT
  chmod +x $LINT

  echo "$(basename $LINT) v$VERSION is now set up."
fi

${LINT} run
