/*
Copyright 2023 The Kubernetes Authors.

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

package common

const (
	// busyboxImage is the image used to extract the e2e logs
	busyboxImage = "registry.k8s.io/e2e-test-images/busybox:1.36.1-1"
	// namespace is the namespace where the conformance pod is created
	namespace = "conformance"
	// PodName is the name of the conformance pod
	PodName = "e2e-conformance-test"
	// ClusterRoleBindingName is the name of the cluster role binding
	ClusterRoleBindingName = "conformance-serviceaccount-role"
	// ClusterRoleName is the name of the cluster role
	ClusterRoleName = "conformance-serviceaccount"
	// ServiceAccountName is the name of the service account
	ServiceAccountName = "conformance-serviceaccount"
	// ConformanceContainer is the name of the conformance container
	ConformanceContainer = "conformance-container"
	// OutputContainer is the name of the busybox container
	OutputContainer = "output-container"
)
