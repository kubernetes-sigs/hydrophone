package cleanup

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/hydrophone/pkg/client"
	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/service"
)

var CleanupCommand = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleanup the resources",
	Long:  `This command will cleanup all the resources.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := client.NewClient()
		config, clientSet := service.Init(viper.GetString("kubeconfig"))
		client.ClientSet = clientSet
		common.PrintInfo(client.ClientSet, config)
		service.Cleanup(client.ClientSet)

		log.Println("Exiting with code: ", client.ExitCode)
		os.Exit(client.ExitCode)
	},
}
