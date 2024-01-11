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

package common

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/kubernetes-sigs/hydrophone/pkg/log"
)

// ArgConfig stores the argument passed when running the program
type ArgConfig struct {
	// Focus set the E2E_FOCUS env var to run a specific test
	// e.g. - sig-auth, sig-apps
	Focus string

	// Skip set the E2E_SKIP env var to skip specified tests
	Skip string

	// ConformanceImage let's people use the conformance container image of their own choice
	// Get the list of images from https://console.cloud.google.com/gcr/images/k8s-artifacts-prod/us/conformance
	// default registry.k8s.io/conformance:v1.29.0
	ConformanceImage string

	// BusyboxImage lets folks use an appropriate busybox image from their own registry
	BusyboxImage string

	// Kubeconfig is the path to the kubeconfig file
	Kubeconfig string

	// Parallel sets the E2E_PARALLEL env var for tests
	Parallel int

	// Verbosity sets the E2E_VERBOSITY env var for tests
	Verbosity int

	// OutputDir is where the e2e.log and junit_01.xml is saved
	OutputDir string

	// Conformance to indicate whether we want to run conformance tests
	ConformanceTests bool

	// DryRun to indicate whether to run ginkgo in dry run mode
	DryRun bool

	// Cleanup indicates we should just cleanup the resources
	Cleanup bool

	// TestRepoList points to the file that has mapping of repositories for test images (KUBE_TEST_REPO_LIST)
	TestRepoList string

	// TestRepo is an alternate repository for test images (KUBE_TEST_REPO)
	TestRepo string

	// ListImages lists images that will be used for conformance tests
	ListImages bool
}

func InitArgs() (*ArgConfig, error) {
	var cfg ArgConfig

	outputDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	flag.StringVar(&cfg.Focus, "focus", "", "focus runs a specific e2e test. e.g. - sig-auth. allows regular expressions.")
	flag.StringVar(&cfg.Skip, "skip", "", "skip specific tests. allows regular expressions.")
	flag.StringVar(&cfg.ConformanceImage, "conformance-image", containerImage,
		"specify a conformance container image of your choice.")
	flag.StringVar(&cfg.BusyboxImage, "busybox-image", busyboxImage,
		"specify an alternate busybox container image.")
	flag.StringVar(&cfg.Kubeconfig, "kubeconfig", "", "path to the kubeconfig file.")
	flag.IntVar(&cfg.Parallel, "parallel", 1, "number of parallel threads in test framework.")
	flag.IntVar(&cfg.Verbosity, "verbosity", 4, "verbosity of test framework.")
	flag.StringVar(&cfg.OutputDir, "output-dir", outputDir, "directory for logs.")
	flag.BoolVar(&cfg.ConformanceTests, "conformance", false, "run conformance tests.")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "run in dry run mode.")
	flag.BoolVar(&cfg.Cleanup, "cleanup", false, "cleanup resources (pods, namespaces etc).")
	flag.StringVar(&cfg.TestRepoList, "test-repo-list", "", "yaml file to override registries for test images.")
	flag.StringVar(&cfg.TestRepo, "test-repo", "", "alternate registry for test images.")
	flag.BoolVar(&cfg.ListImages, "list-images", false, "list all images that will be used during conformance tests.")

	flag.Parse()

	conformance := false
	focus := false
	for _, keyValue := range os.Args {
		arg := strings.Split(keyValue, "=")[0]
		if arg == "--focus" {
			focus = true
		}
		if arg == "--conformance" {
			conformance = true
		}
	}
	if conformance && focus {
		return nil, fmt.Errorf("specify either --conformance or --focus arguments, not both")
	}

	return &cfg, nil
}

func PrintInfo(clientSet *kubernetes.Clientset, config *rest.Config) {
	serverVersion, err := clientSet.ServerVersion()
	if err != nil {
		log.Fatal("Error fetching server version: ", err)
	}

	log.Printf("API endpoint : %s", config.Host)
	log.Printf("Server version : %#v", *serverVersion)
}

func ValidateArgs(clientSet *kubernetes.Clientset, config *rest.Config, cfg *ArgConfig) {
	if cfg.ConformanceTests {
		cfg.Focus = "\\[Conformance\\]"
	}

	if cfg.Skip != "" {
		log.Printf("Skipping tests : '%s'", cfg.Skip)
	}
	log.Printf("Using conformance image : '%s'", cfg.ConformanceImage)
	log.Printf("Using busybox image : '%s'", cfg.BusyboxImage)
	log.Printf("Test framework will start '%d' threads and use verbosity '%d'",
		cfg.Parallel, cfg.Verbosity)

	if _, err := os.Stat(cfg.OutputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(cfg.OutputDir, 0755); err != nil {
			log.Fatalf("Error creating output directory [%s] : %v", cfg.OutputDir, err)
		}
	}
}
