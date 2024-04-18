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

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeVersion(t *testing.T) {
	testCases := []struct {
		name            string
		version         string
		expectedVersion string
		expectErr       bool
	}{
		{
			name:            "stable version",
			version:         "v1.28.6",
			expectedVersion: "v1.28.6",
		},
		{
			name:            "pre released version",
			version:         "v1.28.6+0fb426",
			expectedVersion: "v1.28.6",
		},
		{
			name:            "pre released version with build metadata",
			version:         "v1.28.6+0fb426.20220304",
			expectedVersion: "v1.28.6",
		},
		{
			name:            "invalid version",
			version:         "v1.28,0",
			expectedVersion: "",
			expectErr:       true,
		},
		{
			name:            "short version",
			version:         "v1.28",
			expectedVersion: "",
			expectErr:       true,
		},
		{
			name:            "no v prefix",
			version:         "1.28.6",
			expectedVersion: "v1.28.6",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trimmedVersion, err := normalizeVersion(tc.version)
			assert.Equal(t, tc.expectedVersion, trimmedVersion)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
