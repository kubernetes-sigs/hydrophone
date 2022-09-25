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
}

func InitArgs() ArgConfig {
	var cfg ArgConfig

	flag.StringVar(&cfg.Focus, "focus", "Conformance", "focus runs a specific e2e test. e.g. - sig-auth")
	flag.StringVar(&cfg.Output, "output", "pod_logs", "output lets people get the logs of the pod in a directory")

	flag.Parse()

	return cfg
}
