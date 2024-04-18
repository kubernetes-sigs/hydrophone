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

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateExtraArgs(t *testing.T) {
	testCases := []struct {
		name        string
		extraArgs   []string
		expectedErr string
	}{
		{
			name:        "valid args",
			extraArgs:   []string{"--key1=value1", "--key2=value2"},
			expectedErr: "",
		},
		{
			name:        "invalid: not 2 parts",
			extraArgs:   []string{"invalid-arg"},
			expectedErr: "invalid --extra-args: expected [invalid-arg] to be of --key=value format",
		},
		{
			name:        "invalid: no value",
			extraArgs:   []string{"--key1=value1", "--key2"},
			expectedErr: "invalid --extra-args: expected [--key2] to be of --key=value format",
		},
		{
			name:        "invalid: not a valid flag",
			extraArgs:   []string{"key1=value1", "--key2=value2"},
			expectedErr: "invalid --extra-args: expected key [key1] in [key1=value1] to start with prefix --",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := Configuration{
				ExtraArgs: tc.extraArgs,
			}

			err := config.Validate()
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
