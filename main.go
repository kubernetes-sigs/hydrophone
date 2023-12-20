package main

import (
	"log"

	"github.com/dims/hydrophone/pkg/client"
	"github.com/dims/hydrophone/pkg/service"
)

func main() {
	client := client.NewClient()
	config, clientSet := service.Init()
	client.ClientSet = clientSet

	cfg, err := service.InitArgs()
	if err != nil {
		log.Fatal("Error parsing arguments: ", err)
	}
	serverVersion, err := client.ClientSet.ServerVersion()
	if err != nil {
		log.Fatal("Error fetching server version: ", err)
	}
	log.Printf("API endpoint : %s", config.Host)
	log.Printf("Server version : %#v", *serverVersion)
	log.Printf("Running tests : '%s'", cfg.Focus)
	log.Printf("Using image : '%s'", cfg.Image)

	service.RunE2E(client.ClientSet, cfg)
	client.PrintE2ELogs()
	service.Cleanup(client.ClientSet)
}
