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
	actual := getKubeConfig(kubeconfig)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}

	// Test case 2: kubeconfig is set through environment variable
	kubeconfig = "/path/to/kubeconfig"
	os.Setenv("KUBECONFIG", kubeconfig)
	expected = kubeconfig
	actual = getKubeConfig("")
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}

	// Test case 3: kubeconfig starts with "~"
	kubeconfig = "~/custom/kubeconfig"
	expected = filepath.Join(os.Getenv("HOME"), "custom/kubeconfig")
	actual = getKubeConfig(kubeconfig)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}
