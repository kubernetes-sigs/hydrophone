package service

import (
	"flag"
	"fmt"
)

// ArgConfig stores the argument passed when running the program
type ArgConfig struct {
	// Focus set the E2E_FOCUS env var to run a specific test
	// e.g. - sig-auth, sig-apps
	Focus string

	// Image let's people use the conformance container image of their own choice
	// Get the list of images from https://console.cloud.google.com/gcr/images/k8s-artifacts-prod/us/conformance
	// default registry.k8s.io/conformance:v1.28.0
	Image string
}

func InitArgs() (*ArgConfig, error) {
	var cfg ArgConfig

	flag.StringVar(&cfg.Focus, "focus", "", "focus runs a specific e2e test. e.g. - sig-auth")
	flag.StringVar(&cfg.Image, "image", containerImage,
		"image let's you select your conformance container image of your choice. for example, for v1.25.0 version of tests, use - 'registry.k8s.io/conformance-amd64:v1.25.0'")

	flag.Parse()

	if cfg.Focus == "" {
		return nil, fmt.Errorf("missing --focus argument (use '[Conformance]' to run all conformance tests)")
	}

	return &cfg, nil
}
