package run

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/hydrophone/pkg/client"
	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/service"
)

var (
	focus            string
	skip             string
	conformanceImage string
	busyboxImage     string
	dryRun           bool
	testRepoList     string
	testRepo         string
)

// ConformanceImage is used to define the container image being used to run the conformance tests
var defaultConfImg string

var RunCommand = &cobra.Command{
	Use:   "run",
	Short: "Run conformance tests.",
	Long:  `This command will run all the conformance tests.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := client.NewClient()
		config, clientSet := service.Init(viper.GetString("kubeconfig"))
		client.ClientSet = clientSet

		common.PrintInfo(client.ClientSet, config)

		common.ValidateArgs(client.ClientSet, config)

		service.RunE2E(client.ClientSet)
		client.PrintE2ELogs()
		client.FetchFiles(config, clientSet, viper.GetString("output-dir"))
		client.FetchExitCode()
		service.Cleanup(client.ClientSet)

		log.Println("Exiting with code: ", client.ExitCode)
		os.Exit(client.ExitCode)
	},
}

func init() {
	RunCommand.Flags().StringVar(&focus, "focus", "", "focus runs a specific e2e test. e.g. - sig-auth. allows regular expressions.")
	viper.BindPFlag("focus", RunCommand.Flags().Lookup("focus"))

	RunCommand.Flags().StringVar(&skip, "skip", "", "skip specific tests. allows regular expressions.")
	viper.BindPFlag("skip", RunCommand.Flags().Lookup("skip"))

	RunCommand.Flags().StringVar(&conformanceImage, "conformance-image", "", "specify a conformance container image of your choice.")
	viper.BindPFlag("conformance-image", RunCommand.Flags().Lookup("conformance-image"))

	RunCommand.Flags().StringVar(&busyboxImage, "busybox-image", "", "specify an alternate busybox container image.")
	viper.BindPFlag("busybox-image", RunCommand.Flags().Lookup("busybox-image"))

	RunCommand.Flags().BoolVar(&dryRun, "dry-run", false, "run in dry run mode.")
	viper.BindPFlag("dry-run", RunCommand.Flags().Lookup("dry-run"))

	RunCommand.Flags().StringVar(&testRepoList, "test-repo-list", "", "yaml file to override registries for test images.")
	viper.BindPFlag("test-repo-list", RunCommand.Flags().Lookup("test-repo-list"))

	RunCommand.Flags().StringVar(&testRepo, "test-repo", "", "skip specific tests. allows regular expressions.")
	viper.BindPFlag("test-repo", RunCommand.Flags().Lookup("test-repo"))
}
