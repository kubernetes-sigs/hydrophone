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

package types

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// DefaultBusyboxImage is the image used to extract the e2e logs.
	DefaultBusyboxImage = "registry.k8s.io/e2e-test-images/busybox:1.36.1-1"

	// DefaultNamespace is the default namespace where the conformance pod is created.
	DefaultNamespace = "conformance"
)

type Configuration struct {
	configFile string

	Kubeconfig             string        `yaml:"kubeconfig"`
	Parallel               int           `yaml:"parallel"`
	Verbosity              int           `yaml:"verbosity"`
	OutputDir              string        `yaml:"outputDir"`
	Skip                   string        `yaml:"skip"`
	ConformanceImage       string        `yaml:"conformanceImage"`
	BusyboxImage           string        `yaml:"busyboxImage"`
	Namespace              string        `yaml:"namespace"`
	DryRun                 bool          `yaml:"dryRun"`
	TestRepoList           string        `yaml:"testRepoList"`
	TestRepo               string        `yaml:"testRepo"`
	ExtraArgs              []string      `yaml:"extraArgs"`
	ExtraGinkgoArgs        []string      `yaml:"extraGinkgoArgs"`
	StartupTimeout         time.Duration `yaml:"startupTimeout"`
	DisableProgressStatus  bool          `yaml:"disableProgressStatus"`
	ProgressStatusInterval time.Duration `yaml:"progressStatusInterval"`
}

func NewDefaultConfiguration() Configuration {
	return Configuration{
		Parallel:               1,
		Verbosity:              4,
		OutputDir:              ".",
		BusyboxImage:           DefaultBusyboxImage,
		Namespace:              DefaultNamespace,
		StartupTimeout:         5 * time.Minute,
		DisableProgressStatus:  false,
		ProgressStatusInterval: 30 * time.Second,
	}
}

func (c *Configuration) Validate() error {
	if err := validateArgsFlag(c.ExtraArgs); err != nil {
		return fmt.Errorf("invalid --extra-args: %w", err)
	}

	if err := validateArgsFlag(c.ExtraGinkgoArgs); err != nil {
		return fmt.Errorf("invalid --extra-ginkgo-args: %w", err)
	}

	if c.Parallel > 1 {
		for _, arg := range c.ExtraGinkgoArgs {
			if strings.Contains(arg, "--nodes=") || strings.Contains(arg, "--procs=") {
				return errors.New("--nodes/--procs is automatically set when --parallel is greater than 1")
			}
		}
	}

	return nil
}

func validateArgsFlag(extraArgs []string) error {
	for _, kv := range extraArgs {
		keyValuePair := strings.SplitN(kv, "=", 2)
		if len(keyValuePair) != 2 {
			return fmt.Errorf("expected [%s] to be of --key=value format", kv)
		}

		key := keyValuePair[0]
		if !strings.HasPrefix(key, "--") && strings.Count(key, "--") != 1 {
			return fmt.Errorf("expected key [%s] in [%s] to start with prefix --", key, kv)
		}
	}

	return nil
}

func loadConfiguration(filename string) (*Configuration, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	config := &Configuration{}
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, fmt.Errorf("invalid configuration file: %w", err)
	}

	return config, nil
}
