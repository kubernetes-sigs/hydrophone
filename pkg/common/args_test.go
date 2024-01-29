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

package common

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestValidateArgs(t *testing.T) {
	// Create a temporary output directory
	tempDir := t.TempDir()
	viper.Set("output-dir", tempDir)

	// Set up the test cases
	testCases := []struct {
		name          string
		focus         string
		expectedFocus string
		skip          string
		expectedSkip  string
		extraArgs     []string
		expectedArgs  string
		wantErr       bool
		expectedErr   string
	}{
		{
			name:          "With focus",
			focus:         "\\[E2E\\]",
			expectedFocus: "\\[E2E\\]",
			skip:          "",
			expectedSkip:  "",
			extraArgs:     []string{},
			expectedArgs:  "",
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name:          "With skip",
			focus:         "",
			expectedFocus: "\\[Conformance\\]",
			skip:          "some tests",
			expectedSkip:  "some tests",
			extraArgs:     []string{},
			expectedArgs:  "",
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name:          "With extra args",
			focus:         "",
			expectedFocus: "\\[Conformance\\]",
			skip:          "",
			expectedSkip:  "",
			extraArgs:     []string{"--key1=value1", "--key2=value2"},
			expectedArgs:  "--key1=value1 --key2=value2",
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name:          "Invalid extra args format",
			focus:         "",
			expectedFocus: "\\[Conformance\\]",
			skip:          "",
			expectedSkip:  "",
			extraArgs:     []string{"invalid-arg"},
			expectedArgs:  "",
			wantErr:       true,
			expectedErr:   "expected [[invalid-arg]] in [[invalid-arg]] to be of --key=value format",
		},
		{
			name:          "Extra args with missing values",
			focus:         "",
			expectedFocus: "\\[Conformance\\]",
			skip:          "",
			expectedSkip:  "",
			extraArgs:     []string{"--key1=value1", "--key2"},
			expectedArgs:  "",
			wantErr:       true,
			expectedErr:   "expected [[--key2]] in [[--key1=value1 --key2]] to be of --key=value format",
		},
		{
			name:          "Extra args with invalid key format",
			focus:         "",
			expectedFocus: "\\[Conformance\\]",
			skip:          "",
			expectedSkip:  "",
			extraArgs:     []string{"key1=value1", "--key2=value2"},
			expectedArgs:  "",
			wantErr:       true,
			expectedErr:   "expected key [key1] in [[key1=value1 --key2=value2]] to start with prefix --",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the test environment
			viper.Set("focus", tc.focus)
			viper.Set("skip", tc.skip)
			viper.Set("extra-args", tc.extraArgs)

			// Call the function under test
			err := ValidateArgs()
			if tc.wantErr {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.Nil(t, err)
				if viper.GetString("skip") != tc.expectedSkip {
					t.Errorf("expected skip to be [%s], got [%s]", tc.expectedSkip, viper.GetString("skip"))
				}
				if viper.GetString("focus") != tc.expectedFocus {
					t.Errorf("expected focus to be [%s], got [%s]", tc.expectedFocus, viper.GetString("focus"))
				}
				if viper.GetString("extra-args") != tc.expectedArgs {
					t.Errorf("expected extra-args to be [%s], got [%s]", tc.expectedArgs, viper.GetString("extra-args"))
				}
			}
		})
	}
}
