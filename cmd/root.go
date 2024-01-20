package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/hydrophone/cmd/cleanup"
	"sigs.k8s.io/hydrophone/cmd/listimages"
	"sigs.k8s.io/hydrophone/cmd/run"
	"sigs.k8s.io/hydrophone/pkg/client"
	"sigs.k8s.io/hydrophone/pkg/service"
)

var (
	cfgFile    string
	kubeconfig string
	parallel   int
	verbosity  int
	outputDir  string
	Config     *rest.Config
	Client     *client.Client
	clientSet  *kubernetes.Clientset
)

var rootCmd = &cobra.Command{
	Use:   "hydrohpone",
	Short: "Hydrophone is a lightweight runner for kubernetes tests.",
	Long:  `Hydrophone is a lightweight runner for kubernetes tests.`,
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("Default config file (%s/hydrophone/hydrophone.yaml)", xdg.ConfigHome))

	rootCmd.AddCommand(run.RunCommand)
	rootCmd.AddCommand(cleanup.CleanupCommand)
	rootCmd.AddCommand(listimages.ListImagesCommand)

	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to the kubeconfig file.")

	rootCmd.PersistentFlags().IntVar(&parallel, "parallel", 1, "number of parallel threads in test framework.")
	viper.BindPFlag("parallel", rootCmd.PersistentFlags().Lookup("parallel"))

	rootCmd.PersistentFlags().IntVar(&verbosity, "verbosity", 4, "verbosity of test framework.")
	viper.BindPFlag("verbosity", rootCmd.PersistentFlags().Lookup("verbosity"))

	rootCmd.PersistentFlags().StringVar(&outputDir, "output-dir", workingDir, "directory for logs.")
	viper.BindPFlag("output-dir", rootCmd.PersistentFlags().Lookup("output-dir"))

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
