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

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetKubeConfig(t *testing.T) {
	// Test case 1: kubeconfig is empty
	kubeconfig := ""
	expected := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	actual := GetKubeConfig(kubeconfig)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}

	// Test case 2: kubeconfig is set through environment variable
	kubeconfig = "/path/to/kubeconfig"
	os.Setenv("KUBECONFIG", kubeconfig)
	expected = kubeconfig
	actual = GetKubeConfig("")
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}

	// Test case 3: kubeconfig starts with "~"
	kubeconfig = "~/custom/kubeconfig"
	expected = filepath.Join(os.Getenv("HOME"), "custom/kubeconfig")
	actual = GetKubeConfig(kubeconfig)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}
