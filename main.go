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

package main

import (
	"fmt"
	"os"

	"sigs.k8s.io/hydrophone/pkg/client"
	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/log"
	"sigs.k8s.io/hydrophone/pkg/service"
)

func main() {
	fmt.Println("TEST")
	client := client.NewClient()

	cfg, err := common.InitArgs()
	if err != nil {
		log.Fatal("Error parsing arguments: ", err)
	}

	config, clientSet := service.Init(cfg)
	client.ClientSet = clientSet

	common.PrintInfo(client.ClientSet, config)

	if cfg.ListImages {
		service.PrintListImages(cfg, client.ClientSet, config)
	} else if cfg.Cleanup {
		service.Cleanup(client.ClientSet)
	} else {
		common.ValidateArgs(client.ClientSet, config, cfg)

		service.RunE2E(client.ClientSet, cfg)
		client.PrintE2ELogs()
		client.FetchFiles(config, clientSet, cfg.OutputDir)
		client.FetchExitCode()
		service.Cleanup(client.ClientSet)
	}
	log.Println("Exiting with code: ", client.ExitCode)
	os.Exit(client.ExitCode)
}
