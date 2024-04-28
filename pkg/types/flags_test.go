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
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveKubeconfig(t *testing.T) {
	// reset environment
	os.Setenv("KUBECONFIG", "")

	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	// Test case 1: kubeconfig is empty
	kubeconfig := ""
	expected := filepath.Join(homeDir, ".kube", "config")

	actual, err := resolveKubeconfig(kubeconfig)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Test case 2: kubeconfig is set through environment variable
	kubeconfig = "/path/to/kubeconfig"
	os.Setenv("KUBECONFIG", kubeconfig)
	expected = kubeconfig

	actual, err = resolveKubeconfig("")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Test case 3: kubeconfig contains with "~"
	kubeconfig = "~/custom/kubeconfig"
	expected = filepath.Join(homeDir, "custom", "kubeconfig")

	actual, err = resolveKubeconfig(kubeconfig)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestMergeConfigs(t *testing.T) {
	testcases := []struct {
		name          string
		flagConfig    Configuration
		loaded        Configuration
		changedFields []string
		expected      Configuration
	}{
		{
			name:          "empty testcase",
			flagConfig:    Configuration{},
			loaded:        Configuration{},
			changedFields: nil,
			expected:      Configuration{},
		},
		{
			name:          "no config, only --output-dir given",
			flagConfig:    Configuration{OutputDir: "foo"},
			loaded:        Configuration{},
			changedFields: []string{"output-dir"},
			expected:      Configuration{OutputDir: "foo"},
		},
		{
			name:          "only load config from file",
			flagConfig:    Configuration{},
			loaded:        Configuration{OutputDir: "foo"},
			changedFields: nil,
			expected:      Configuration{OutputDir: "foo"},
		},
		{
			name:          "CLI flag has priority over loaded config file",
			flagConfig:    Configuration{OutputDir: "foo"},
			loaded:        Configuration{OutputDir: "bar"},
			changedFields: []string{"output-dir"},
			expected:      Configuration{OutputDir: "foo"},
		},
		{
			name:          "CLI flag has priority, but only if a non-empty value is given",
			flagConfig:    Configuration{OutputDir: ""},
			loaded:        Configuration{OutputDir: "bar"},
			changedFields: []string{"output-dir"},
			expected:      Configuration{OutputDir: "bar"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			merged := mergeConfigs(func(flag string) bool {
				return slices.Contains(tc.changedFields, flag)
			}, &tc.flagConfig, &tc.loaded)

			assert.NotNil(t, merged)
			assert.Equal(t, tc.expected, *merged)
		})
	}
}
