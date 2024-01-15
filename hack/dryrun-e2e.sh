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

set -o errexit -o nounset -o xtrace

# Default versions k8s and kind
K8S_VERSION=${K8S_VERSION:-v1.29.0}
KIND_VERSION=${KIND_VERSION:-v0.20.0}

# Maximum time (in seconds) for a dry run test
DRYRUN_THRESHOLD=${DRYRUN_DURATION:-5}

# Set the artifacts directory, defaulting to a local subdirectory
export ARTIFACTS="${ARTIFACTS:-${PWD}/_artifacts}"
mkdir -p "${ARTIFACTS}/results"

# Download and install kind
curl -fsSL -o ./kind "https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-linux-amd64"
install --mode=755 ./kind /usr/local/bin/kind

# Create a kind cluster with a specific Kubernetes version that will match the hydrophone test
cat <<EOF | kind create cluster --image kindest/node:${K8S_VERSION} --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
EOF

# Retrieve cluster information
kubectl cluster-info --context kind-kind
kubectl get nodes

# Execute hydrophone with specific parameters and log output
bin/hydrophone \
  --focus 'Simple pod should contain last line of the log' \
  --output-dir ${ARTIFACTS}/results/ \
  --conformance-image registry.k8s.io/conformance:${K8S_VERSION} \
  --dry-run | tee /tmp/dryrun.log

# Check the duration of the dry run against the threshold
DRYRUN_DURATION=$(grep -oP 'Ran 1 of \d+ Specs in \K[0-9.]+(?= seconds)' /tmp/dryrun.log | cut -d. -f1)
if [[ ${DRYRUN_DURATION} -gt ${DRYRUN_THRESHOLD} ]]; then 
  echo "Focused test took too long to run. Expected less than ${DRYRUN_THRESHOLD} seconds, got ${DRYRUN_DURATION} seconds"
  exit 1
fi