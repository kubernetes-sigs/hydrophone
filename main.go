package main

import (
	"github.com/dims/k8s-run-e2e/pkg/client"
	"github.com/dims/k8s-run-e2e/pkg/service"
)

func main() {
	c := client.NewClient()
	c.ClientSet = service.Init()

	
}
