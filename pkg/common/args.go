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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/hydrophone/pkg/log"
)

// SetDefaults sets the default values for various configuration options used in the application.
// Finally, it logs the API endpoint, server version, namespace, conformance image, and busybox image.
func SetDefaults(clientset *kubernetes.Clientset, config *rest.Config) error {
	serverVersion, err := clientset.ServerVersion()
	if err != nil {
		return fmt.Errorf("failed fetching server version: %w", err)
	}
	trimmedVersion, err := trimVersion(serverVersion.String())
	if err != nil {
		return fmt.Errorf("failed parsing server version: %w", err)
	}

	if viper.Get("conformance-image") == "" {
		viper.Set("conformance-image", fmt.Sprintf("registry.k8s.io/conformance:%s", trimmedVersion))
	}
	if viper.Get("busybox-image") == "" {
		viper.Set("busybox-image", busyboxImage)
	}
	if viper.Get("namespace") == "" {
		viper.Set("namespace", DefaultNamespace)
	}

	log.Printf("API endpoint: %s", config.Host)
	log.Printf("Server version: %#v", *serverVersion)
	log.Printf("Using namespace: %s", viper.Get("namespace"))
	log.Printf("Using conformance image: %s", viper.Get("conformance-image"))
	log.Printf("Using busybox image: %s", viper.Get("busybox-image"))

	if viper.Get("skip") != "" {
		log.Printf("Skipping tests: %s", viper.Get("skip"))
	}

	if viper.GetBool("dry-run") {
		log.Println("Dry-run enabled.")
	}

	log.Printf("Test framework will start %d thread(s) and use verbosity level %d.",
		viper.Get("parallel"), viper.Get("verbosity"))

	return nil
}

// ValidateConformanceArgs validates the arguments passed to the program
// and creates the output directory if it doesn't exist
func ValidateConformanceArgs() error {
	if viper.Get("focus") == "" {
		viper.Set("focus", "\\[Conformance\\]")
	}

	if err := validateArgsFlag(viper.GetStringSlice("extra-args")); err != nil {
		return fmt.Errorf("invalid --extra-args: %w", err)
	}

	if err := validateArgsFlag(viper.GetStringSlice("extra-ginkgo-args")); err != nil {
		return fmt.Errorf("invalid --extra-ginkgo-args: %w", err)
	}

	if viper.GetInt("parallel") > 1 {
		extraArgs := viper.GetStringSlice("extra-ginkgo-args")
		for _, arg := range extraArgs {
			if strings.Contains(arg, "--nodes=") || strings.Contains(arg, "--procs=") {
				return errors.New("--nodes/--procs is automatically set when --parallel is greater than 1")
			}
		}
	}

	outputDir := viper.GetString("output-dir")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("error creating output directory [%s]: %w", outputDir, err)
	}

	return nil
}

func validateArgsFlag(extraArgs []string) error {
	if len(extraArgs) != 0 {
		for _, kv := range extraArgs {
			keyValuePair := strings.SplitN(kv, "=", 2)
			if len(keyValuePair) != 2 {
				return fmt.Errorf("expected [%s] in [%s] to be of --key=value format", keyValuePair, extraArgs)
			}
			key := keyValuePair[0]
			if !strings.HasPrefix(key, "--") && strings.Count(key, "--") != 1 {
				return fmt.Errorf("expected key [%s] in [%s] to start with prefix --", key, extraArgs)
			}
		}
	}

	return nil
}

func trimVersion(version string) (string, error) {
	version = strings.TrimPrefix(version, "v")

	parsedVersion, err := semver.Parse(version)
	if err != nil {
		return "", fmt.Errorf("error parsing conformance image tag: %w", err)
	}

	return "v" + parsedVersion.FinalizeVersion(), nil
}
