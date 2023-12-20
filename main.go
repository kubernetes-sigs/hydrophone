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
	"log"

	"github.com/dims/hydrophone/pkg/client"
	"github.com/dims/hydrophone/pkg/service"
)

func main() {
	client := client.NewClient()

	cfg, err := service.InitArgs()
	if err != nil {
		log.Fatal("Error parsing arguments: ", err)
	}

	config, clientSet := service.Init(cfg)
	client.ClientSet = clientSet

	serverVersion, err := client.ClientSet.ServerVersion()
	if err != nil {
		log.Fatal("Error fetching server version: ", err)
	}
	log.Printf("API endpoint : %s", config.Host)
	log.Printf("Server version : %#v", *serverVersion)
	log.Printf("Running tests : '%s'", cfg.Focus)
	if cfg.Skip != "" {
		log.Printf("Skipping tests : '%s'", cfg.Skip)
	}
	log.Printf("Using image : '%s'", cfg.Image)

	service.RunE2E(client.ClientSet, cfg)
	client.PrintE2ELogs()
	service.Cleanup(client.ClientSet)
}
