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

//nolint:errcheck // TODO(embik): temporary put in place to ensure PRs do not conflict
package cmd

import (
	"fmt"
	"os"

	"sigs.k8s.io/hydrophone/pkg/client"
	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/log"
	"sigs.k8s.io/hydrophone/pkg/service"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile          string
	kubeconfig       string
	parallel         int
	verbosity        int
	outputDir        string
	cleanup          bool
	listImages       bool
	conformance      bool
	focus            string
	skip             string
	conformanceImage string
	busyboxImage     string
	namespace        string
	dryRun           bool
	testRepoList     string
	testRepo         string
)

var rootCmd = &cobra.Command{
	Use:   "hydrophone",
	Short: "Hydrophone is a lightweight runner for Kubernetes tests.",
	Long:  `Hydrophone is a lightweight runner for Kubernetes tests.`,
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()

		client := client.NewClient()
		config, clientSet, err := service.Init(viper.GetString("kubeconfig"))
		if err != nil {
			log.Fatalf("Failed to init kube client: %v.", err)
		}

		client.ClientSet = clientSet
		err = common.SetDefaults(client.ClientSet, config)
		if err != nil {
			log.Fatalf("Failed to apply default values: %v.", err)
		}

		if cleanup {
			if err := service.Cleanup(ctx, client.ClientSet); err != nil {
				log.Fatalf("Failed to cleanup: %v.", err)
			}
		} else if listImages {
			if err := service.PrintListImages(ctx, client.ClientSet); err != nil {
				log.Fatalf("Failed to list images: %v.", err)
			}
		} else {
			if err := common.ValidateConformanceArgs(); err != nil {
				log.Fatalf("Invalid arguments: %v.", err)
			}

			verboseGinkgo := verbosity >= 6
			showSpinner := !verboseGinkgo && verbosity > 2

			if err := service.RunE2E(ctx, client.ClientSet, verboseGinkgo); err != nil {
				log.Fatalf("Failed to run tests: %v.", err)
			}

			var spinner *common.Spinner
			if showSpinner {
				spinner = common.NewSpinner(os.Stdout)
				spinner.Start()
			}

			// PrintE2ELogs is a long running method
			if err := client.PrintE2ELogs(ctx); err != nil {
				log.Fatalf("Failed to get test logs: %v.", err)
			}

			if showSpinner {
				spinner.Stop()
			}

			if err := client.FetchFiles(ctx, config, clientSet, viper.GetString("output-dir")); err != nil {
				log.Fatalf("Failed to download results: %v.", err)
			}
			if err := client.FetchExitCode(ctx); err != nil {
				log.Fatalf("Failed to determine exit code: %v.", err)
			}
			if err := service.Cleanup(ctx, client.ClientSet); err != nil {
				log.Fatalf("Failed to cleanup: %v.", err)
			}

			if client.ExitCode == 0 {
				log.Println("Tests completed successfully.")
			} else {
				log.Printf("Tests failed (code %d).", client.ExitCode)
				os.Exit(client.ExitCode)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (defaults to %s/hydrophone.yaml).", xdg.ConfigHome))
	rootCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "path to the kubeconfig file.")

	rootCmd.Flags().IntVar(&parallel, "parallel", 1, "number of parallel threads in test framework.")
	viper.BindPFlag("parallel", rootCmd.Flags().Lookup("parallel"))

	rootCmd.Flags().IntVar(&verbosity, "verbosity", 4, "verbosity of test framework.")
	viper.BindPFlag("verbosity", rootCmd.Flags().Lookup("verbosity"))

	rootCmd.Flags().StringVar(&outputDir, "output-dir", workingDir, "directory for logs.")
	viper.BindPFlag("output-dir", rootCmd.Flags().Lookup("output-dir"))

	rootCmd.Flags().BoolVar(&cleanup, "cleanup", false, "cleanup resources (pods, namespaces etc).")

	rootCmd.Flags().BoolVar(&listImages, "list-images", false, "list all images that will be used during conformance tests.")

	rootCmd.Flags().BoolVar(&conformance, "conformance", false, "run conformance tests.")

	rootCmd.Flags().StringVar(&focus, "focus", "", "focus runs a specific e2e test. e.g. - sig-auth. allows regular expressions.")
	viper.BindPFlag("focus", rootCmd.Flags().Lookup("focus"))

	rootCmd.Flags().StringVar(&skip, "skip", "", "skip specific tests. allows regular expressions.")
	viper.BindPFlag("skip", rootCmd.Flags().Lookup("skip"))

	rootCmd.Flags().StringVar(&conformanceImage, "conformance-image", "", "specify a conformance container image of your choice.")
	viper.BindPFlag("conformance-image", rootCmd.Flags().Lookup("conformance-image"))

	rootCmd.Flags().StringVar(&busyboxImage, "busybox-image", "", "specify an alternate busybox container image.")
	viper.BindPFlag("busybox-image", rootCmd.Flags().Lookup("busybox-image"))

	rootCmd.Flags().StringVar(&namespace, "namespace", "", "the namespace where the conformance pod is created.")
	viper.BindPFlag("namespace", rootCmd.Flags().Lookup("namespace"))

	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "run in dry run mode.")
	viper.BindPFlag("dry-run", rootCmd.Flags().Lookup("dry-run"))

	rootCmd.Flags().StringVar(&testRepoList, "test-repo-list", "", "yaml file to override registries for test images.")
	viper.BindPFlag("test-repo-list", rootCmd.Flags().Lookup("test-repo-list"))

	rootCmd.Flags().StringVar(&testRepo, "test-repo", "", "skip specific tests. allows regular expressions.")
	viper.BindPFlag("test-repo", rootCmd.Flags().Lookup("test-repo"))

	rootCmd.Flags().StringSlice("extra-args", []string{}, "Additional parameters to be provided to the conformance container. These parameters should be specified as key-value pairs, separated by commas. Each parameter should start with -- (e.g., --clean-start=true,--allowed-not-ready-nodes=2)")
	viper.BindPFlag("extra-args", rootCmd.Flags().Lookup("extra-args"))

	rootCmd.MarkFlagsMutuallyExclusive("conformance", "focus", "cleanup", "list-images")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// the config will be located under `~/.config/hydrophone.yaml` on linux
		configDir := xdg.ConfigHome
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("hydrophone")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				err := viper.SafeWriteConfig()
				if err != nil {
					log.Fatal("Error:", err)
				}
			} else {
				log.Fatal(err)
			}
		}
	}
	kubeconfig = service.GetKubeConfig(kubeconfig)
	viper.Set("kubeconfig", kubeconfig)
}
