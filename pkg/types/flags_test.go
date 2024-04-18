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
