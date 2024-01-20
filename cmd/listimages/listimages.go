package listimages

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/hydrophone/pkg/client"
	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/service"
)

var ListImagesCommand = &cobra.Command{
	Use:   "listimages",
	Short: "List all images",
	Long:  `This command will list all images that will be used for conformance tests.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := client.NewClient()
		config, clientSet := service.Init(viper.GetString("kubeconfig"))
		client.ClientSet = clientSet
		common.PrintInfo(client.ClientSet, config)
		service.PrintListImages(client.ClientSet)

		log.Println("Exiting with code: ", client.ExitCode)
		os.Exit(client.ExitCode)
	},
}
