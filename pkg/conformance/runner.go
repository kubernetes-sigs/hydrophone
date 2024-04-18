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

package conformance

import (
	"fmt"

	"sigs.k8s.io/hydrophone/pkg/types"

	"k8s.io/client-go/kubernetes"
)

type TestRunner struct {
	config    types.Configuration
	clientset *kubernetes.Clientset
}

func NewTestRunner(config types.Configuration, clientset *kubernetes.Clientset) *TestRunner {
	return &TestRunner{
		config:    config,
		clientset: clientset,
	}
}

func (r *TestRunner) namespacedName(basename string) string {
	return fmt.Sprintf("%s:%s", basename, r.config.Namespace)
}
