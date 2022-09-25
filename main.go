package main

import (
	"github.com/dims/k8s-run-e2e/pkg/client"
	"github.com/dims/k8s-run-e2e/pkg/service"
)

func main() {
	client := client.NewClient()
	client.ClientSet = service.Init()

	cfg := service.InitArgs()
	service.RunE2E(client.ClientSet, cfg.Focus)
	client.CheckForE2ELogs(cfg.Output)
	service.Cleanup(client.ClientSet)
}
