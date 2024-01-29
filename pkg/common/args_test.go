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
	"reflect"
	"testing"

	"github.com/spf13/viper"
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
		expectedArgs  []string
	}{
		{
			name:          "With focus",
			focus:         "\\[E2E\\]",
			expectedFocus: "\\[E2E\\]",
			skip:          "",
			expectedSkip:  "",
			extraArgs:     []string{},
			expectedArgs:  []string{},
		},
		{
			name:          "With skip",
			focus:         "",
			expectedFocus: "\\[Conformance\\]",
			skip:          "some tests",
			expectedSkip:  "some tests",
			extraArgs:     []string{},
			expectedArgs:  []string{},
		},
		{
			name:          "With extra args",
			focus:         "",
			expectedFocus: "\\[Conformance\\]",
			skip:          "",
			expectedSkip:  "",
			extraArgs:     []string{"--key1=value1", "--key2=value2"},
			expectedArgs:  []string{"--key1=value1", "--key2=value2"},
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
			ValidateArgs()
			if viper.GetString("skip") != tc.expectedSkip {
				t.Errorf("expected skip to be [%s], got [%s]", tc.expectedSkip, viper.GetString("skip"))
			}
			if viper.GetString("focus") != tc.expectedFocus {
				t.Errorf("expected focus to be [%s], got [%s]", tc.expectedFocus, viper.GetString("focus"))
			}
			if !reflect.DeepEqual(viper.GetStringSlice("extra-args"), tc.expectedArgs) {
				t.Errorf("expected extra-args to be [%v], got [%v]", tc.expectedArgs, viper.GetStringSlice("extra-args"))
			}
		})
	}
}
