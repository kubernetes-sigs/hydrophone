package main

import (
	"github.com/dims/k8s-run-e2e/pkg/client"
	"github.com/dims/k8s-run-e2e/pkg/service"
	"log"
)

func main() {
	client := client.NewClient()
	client.ClientSet = service.Init()

	cfg := service.InitArgs()

	if cfg.Focus == "" {
		log.Fatal("please specify which tests to run using --focus argument\n" +
			"(for example '[Conformance]' to run all conformance tests)")
	}

	service.RunE2E(client.ClientSet, cfg)
	client.CheckForE2ELogs(cfg.Output)
	service.Cleanup(client.ClientSet)
}
