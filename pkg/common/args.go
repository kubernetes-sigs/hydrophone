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

package common

import (
	"fmt"
	"strings"
	"time"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	
	"sigs.k8s.io/hydrophone/pkg/log"
)

// PrintInfo prints the information about the cluster
func PrintInfo(clientSet *kubernetes.Clientset, config *rest.Config) {
	spinner := NewSpinner(os.Stdout)
    spinner.Start()
	 
	time.Sleep(2 * time.Second)
	serverVersion, err := clientSet.ServerVersion()
	if err != nil {
		log.Fatal("Error fetching server version: ", err)
	}
	if viper.Get("conformance-image") == "" {
		viper.Set("conformance-image", fmt.Sprintf("registry.k8s.io/conformance:%s", serverVersion.String()))
	}
	if viper.Get("busybox-image") == "" {
		viper.Set("busybox-image", busyboxImage)
	}

	log.printf("API endpoint : %s", config.Host)
	log.Printf("Server version : %#v", *serverVersion)
}

// ValidateArgs validates the arguments passed to the program
// and creates the output directory if it doesn't exist

func ValidateArgs() error {
	if viper.Get("namespace") == "" {
		viper.Set("namespace", DefaultNamespace)
	}
	if viper.Get("focus") == "" {
		viper.Set("focus", "\\[Conformance\\]")
	}

	if viper.Get("skip") != "" {
		log.Printf("Skipping tests : '%s'", viper.Get("skip"))
	}

	if extraArgs := viper.GetStringSlice("extra-args"); len(extraArgs) != 0 {
		for _, kv := range extraArgs {
			keyValuePair := strings.SplitN(kv, "=", 2)
			if len(keyValuePair) != 2 {
				return fmt.Errorf("expected [%s] in [%s] to be of --key=value format", keyValuePair, extraArgs)
			}
			key := keyValuePair[0]
			if !strings.HasPrefix(key, "--") && strings.Count(key, "--") != 1 {
				return fmt.Errorf("expected key [%s] in [%s] to start with prefix --", key, extraArgs)
			}
		}
	}

	log.Printf("Using namespace : '%s'", viper.Get("namespace"))
	log.Printf("Using conformance image : '%s'", viper.Get("conformance-image"))
	log.Printf("Using busybox image : '%s'", viper.Get("busybox-image"))
	log.Printf("Test framework will start '%d' threads and use verbosity '%d'",
		viper.Get("parallel"), viper.Get("verbosity"))

	outputDir := viper.GetString("output-dir")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory [%s] : %v", outputDir, err)
		}
	}
	return nil
}