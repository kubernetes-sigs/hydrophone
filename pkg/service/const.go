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

package service

const (
	containerImage         = "registry.k8s.io/conformance:v1.28.0"
	Namespace              = "conformance"
	PodName                = "e2e-conformance-test"
	ClusterRoleBindingName = "conformance-serviceaccount-role"
	ClusterRoleName        = "conformance-serviceaccount"
	ServiceAccountName     = "conformance-serviceaccount"
	ConformanceContainer   = "conformance-container"
	OutputContainer        = "output-container"
)
