/*
Copyright 2024 The Kubernetes Authors.

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
	"testing"

	"sigs.k8s.io/hydrophone/pkg/types"

	"github.com/stretchr/testify/require"
)

func TestNamespacedName(t *testing.T) {
	tests := []struct {
		name      string
		basename  string
		namespace string
		expected  string
	}{
		{
			name:      "basic test",
			basename:  "testrole",
			namespace: "default",
			expected:  "testrole:default",
		},
		{
			name:      "empty namespace",
			basename:  "foo",
			namespace: "",
			expected:  "foo:",
		},
		{
			name:      "empty basename",
			basename:  "",
			namespace: "dev",
			expected:  ":dev",
		},
		{
			name:      "both empty",
			basename:  "",
			namespace: "",
			expected:  ":",
		},
		{
			name:      "special chars",
			basename:  "test-pod-123",
			namespace: "kube-system",
			expected:  "test-pod-123:kube-system",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			runner := NewTestRunner(types.Configuration{
				Namespace: tt.namespace,
			}, nil)

			result := runner.namespacedName(tt.basename)
			require.Equal(t, tt.expected, result)
		})
	}
}
