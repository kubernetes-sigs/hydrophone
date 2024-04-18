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

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/conformance"
	"sigs.k8s.io/hydrophone/pkg/conformance/client"
	"sigs.k8s.io/hydrophone/pkg/log"
	"sigs.k8s.io/hydrophone/pkg/types"

	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	runCleanup       bool
	runListImages    bool
	runConformance   bool
	conformanceFocus string
)

func New() *cobra.Command {
	var (
		rootCmd *cobra.Command
		config  types.Configuration
	)

	rootCmd = &cobra.Command{
		Use:   "hydrophone",
		Short: "Hydrophone is a lightweight runner for Kubernetes tests.",
		Long:  "Hydrophone is a lightweight runner for Kubernetes tests.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			effectiveConfig, err := config.Complete(rootCmd.Flags())
			if err != nil {
				_ = rootCmd.Usage()
				return err
			}

			return action(cmd.Context(), effectiveConfig)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	config = types.NewDefaultConfiguration()
	config.AddFlags(rootCmd.Flags())

	// the different ways to run hydrophone are not part of the configuration file
	rootCmd.Flags().BoolVar(&runCleanup, "cleanup", false, "cleanup resources (pods, namespaces etc).")
	rootCmd.Flags().BoolVar(&runListImages, "list-images", false, "list all images that will be used during conformance tests.")
	rootCmd.Flags().BoolVar(&runConformance, "conformance", false, "run conformance tests.")
	rootCmd.Flags().StringVar(&conformanceFocus, "focus", "", "focus runs a specific e2e test. e.g. - sig-auth. allows regular expressions.")

	rootCmd.MarkFlagsMutuallyExclusive("conformance", "focus", "cleanup", "list-images")

	return rootCmd
}

func action(ctx context.Context, config *types.Configuration) error {
	if err := os.MkdirAll(config.OutputDir, 0o755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	restConfig, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		return fmt.Errorf("error loading kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("error getting config client: %w", err)
	}

	// some defaults can only be applied after we connected to the cluster
	config, serverVersion, err := applyClusterDefaults(config, clientset)
	if err != nil {
		return fmt.Errorf("error applying cluster configuration: %w", err)
	}

	// print effective runtime config before we begin
	log.Printf("API endpoint: %s", restConfig.Host)
	log.Printf("Server version: %#v", *serverVersion)
	log.Printf("Using namespace: %s", config.Namespace)
	log.Printf("Using conformance image: %s", config.ConformanceImage)
	log.Printf("Using busybox image: %s", config.BusyboxImage)

	if config.Skip != "" {
		log.Printf("Skipping tests: %s", config.Skip)
	}

	if config.DryRun {
		log.Println("Dry-run enabled.")
	}

	log.Printf("Test framework will start %d thread(s) and use verbosity level %d.", config.Parallel, config.Verbosity)

	// prepare test runner and the client to monitor it
	testRunner := conformance.NewTestRunner(*config, clientset)
	testClient := client.NewClient(restConfig, clientset, config.Namespace)

	switch {
	case runCleanup:
		if err := testRunner.Cleanup(ctx); err != nil {
			return fmt.Errorf("failed to cleanup: %w", err)
		}

	case runListImages:
		if err := testRunner.PrintListImages(ctx, config.StartupTimeout); err != nil {
			return fmt.Errorf("failed to list images: %w", err)
		}

	default:
		// `hydrophone --conformance` is an alias for `hydrophone --focus '\[Conformance\]'`
		if conformanceFocus == "" {
			conformanceFocus = `\[Conformance\]`
		}

		verboseGinkgo := config.Verbosity >= 6
		showSpinner := !verboseGinkgo && config.Verbosity > 2

		if err := testRunner.Deploy(ctx, conformanceFocus, verboseGinkgo, config.StartupTimeout); err != nil {
			return fmt.Errorf("failed to deploy tests: %w", err)
		}

		before := time.Now()

		var spinner *common.Spinner
		if showSpinner {
			spinner = common.NewSpinner(os.Stdout)
			spinner.Start()
		}

		// PrintE2ELogs is a long running method
		if err := testClient.PrintE2ELogs(ctx); err != nil {
			return fmt.Errorf("failed to get test logs: %w", err)
		}

		if showSpinner {
			spinner.Stop()
		}

		log.Printf("Tests finished after %v.", time.Since(before).Round(time.Second))

		if err := testClient.FetchFiles(ctx, config.OutputDir); err != nil {
			return fmt.Errorf("failed to download results: %w", err)
		}

		exitCode, err := testClient.FetchExitCode(ctx)
		if err != nil {
			return fmt.Errorf("failed to determine exit code: %w", err)
		}

		if err := testRunner.Cleanup(ctx); err != nil {
			return fmt.Errorf("failed to cleanup: %w", err)
		}

		if exitCode == 0 {
			log.Println("Tests completed successfully.")
		} else {
			log.Printf("Tests failed (code %d).", exitCode)
			os.Exit(exitCode)
		}
	}

	return nil
}

func applyClusterDefaults(config *types.Configuration, clientset *kubernetes.Clientset) (*types.Configuration, *version.Info, error) {
	serverVersion, err := clientset.ServerVersion()
	if err != nil {
		return nil, nil, fmt.Errorf("failed fetching server version: %w", err)
	}

	normalized, err := normalizeVersion(serverVersion.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed parsing server version: %w", err)
	}

	if config.ConformanceImage == "" {
		config.ConformanceImage = fmt.Sprintf("registry.k8s.io/conformance:%s", normalized)
	}

	return config, serverVersion, nil
}

func normalizeVersion(version string) (string, error) {
	version = strings.TrimPrefix(version, "v")

	parsedVersion, err := semver.Parse(version)
	if err != nil {
		return "", fmt.Errorf("error parsing conformance image tag: %w", err)
	}

	return "v" + parsedVersion.FinalizeVersion(), nil
}
