package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/hydrophone/cmd/run"
	"sigs.k8s.io/hydrophone/pkg/client"
	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/service"
)

var (
	cfgFile    string
	kubeconfig string
	parallel   int
	verbosity  int
	outputDir  string
	cleanup    bool
	listImages bool
)

var rootCmd = &cobra.Command{
	Use:   "hydrohpone",
	Short: "Hydrophone is a lightweight runner for kubernetes tests.",
	Long:  `Hydrophone is a lightweight runner for kubernetes tests.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cleanup {
			client := client.NewClient()
			config, clientSet := service.Init(viper.GetString("kubeconfig"))
			client.ClientSet = clientSet
			common.PrintInfo(client.ClientSet, config)
			service.Cleanup(client.ClientSet)

			log.Println("Exiting with code: ", client.ExitCode)
			os.Exit(client.ExitCode)
		}
		if listImages {
			client := client.NewClient()
			config, clientSet := service.Init(viper.GetString("kubeconfig"))
			client.ClientSet = clientSet
			common.PrintInfo(client.ClientSet, config)
			service.PrintListImages(client.ClientSet)

			log.Println("Exiting with code: ", client.ExitCode)
			os.Exit(client.ExitCode)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(run.RunCommand)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("Default config file (%s/hydrophone/hydrophone.yaml)", xdg.ConfigHome))
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to the kubeconfig file.")

	rootCmd.PersistentFlags().IntVar(&parallel, "parallel", 1, "number of parallel threads in test framework.")
	viper.BindPFlag("parallel", rootCmd.PersistentFlags().Lookup("parallel"))

	rootCmd.PersistentFlags().IntVar(&verbosity, "verbosity", 4, "verbosity of test framework.")
	viper.BindPFlag("verbosity", rootCmd.PersistentFlags().Lookup("verbosity"))

	rootCmd.PersistentFlags().StringVar(&outputDir, "output-dir", workingDir, "directory for logs.")
	viper.BindPFlag("output-dir", rootCmd.PersistentFlags().Lookup("output-dir"))

	rootCmd.Flags().BoolVar(&cleanup, "cleanup", false, "cleanup resources (pods, namespaces etc).")
	viper.BindPFlag("cleanup", rootCmd.Flags().Lookup("cleanup"))

	rootCmd.Flags().BoolVar(&listImages, "list-images", false, "list all images that will be used during conformance tests.")
	viper.BindPFlag("list-images", rootCmd.Flags().Lookup("list-images"))
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// the config will belocated under `~/.config/hydrophone/hydrophone.yaml` on linux
		configDir := xdg.ConfigHome
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("hydrophone")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				err := viper.SafeWriteConfig()
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println(err)
			}
		}
	}

	if err := viper.ReadInConfig(); err == nil {
		_ = 1
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	kubeconfig = service.GetKubeConfig(kubeconfig)
	viper.Set("kubeconfig", kubeconfig)
}
