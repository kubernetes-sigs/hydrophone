#!/bin/bash

set -xeuo pipefail

HYDROPHONE_ROOT=$(git rev-parse --show-toplevel)
echo "HYDROPHONE_ROOT: $HYDROPHONE_ROOT"

pushd "${HYDROPHONE_ROOT}" >/dev/null
  go mod edit -json | jq -r ".Require[] | .Path | select(contains(\"k8s.io/\"))" | xargs xargs -L1 go get -d
  go mod tidy

  K8S_VERSION=$(curl https://cdn.dl.k8s.io/release/stable.txt -s)
  sed -i "s|K8S_VERSION: .*|K8S_VERSION: $K8S_VERSION|" .github/workflows/*.yml
  sed -i -r "s/conformance:v[0-9]+\.[0-9]+\.[0-9]+/conformance:$K8S_VERSION/g" README.md pkg/common/*.go

popd >/dev/null
git status