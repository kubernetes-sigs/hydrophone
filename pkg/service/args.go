package service

import (
	"flag"
)

// ArgConfig stores the argument passed when running the program
type ArgConfig struct {
	// Output specifies the directory where pod logs will be stored
	Output string

	// Focus set the E2E_FOCUS env var to run a specific test
	// e.g. - sig-auth, sig-apps
	Focus string

	// Image let's people use the conformance container image of their own choice
	// Get the list of images from https://console.cloud.google.com/gcr/images/k8s-artifacts-prod/us/conformance
	Image string
}

func InitArgs() ArgConfig {
	var cfg ArgConfig

	flag.StringVar(&cfg.Focus, "focus", "Conformance", "focus runs a specific e2e test. e.g. - sig-auth")
	flag.StringVar(&cfg.Output, "output", "pod_logs", "output lets people get the logs of the pod in a directory")
	flag.StringVar(&cfg.Image, "image", containerImage, "image let's you select your conformance container image of yoor choice")

	flag.Parse()

	return cfg
}
