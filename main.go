package main

import (
	"github.com/dims/k8s-run-e2e/pkg/client"
	"github.com/dims/k8s-run-e2e/pkg/service"
)

func main() {
	client := client.NewClient()
	client.ClientSet = service.Init()

	client.CheckForE2ELogs()
}
