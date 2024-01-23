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

package client

import (
	"context"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/log"
)

var (
	ctx = context.TODO()
)

// Client is a struct that holds the clientset and exit code
type Client struct {
	ClientSet *kubernetes.Clientset
	ExitCode  int
}

// FetchFiles downloads the e2e.log and junit_01.xml files from the pod
// and writes them to the output directory
func (c *Client) FetchFiles(config *rest.Config, clientset *kubernetes.Clientset, outputDir string) {
	log.Println("downloading e2e.log to ", filepath.Join(outputDir, "e2e.log"))
	e2eLogFile, err := os.OpenFile(filepath.Join(outputDir, "e2e.log"), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("unable to create e2e.log: %v\n", err)
	}
	defer e2eLogFile.Close()
	err = downloadFile(config, clientset, common.PodName, common.OutputContainer, "/tmp/results/e2e.log", e2eLogFile)
	if err != nil {
		log.Fatalf("unable to download e2e.log: %v\n", err)
	}
	log.Println("downloading junit_01.xml to", filepath.Join(outputDir, "junit_01.xml"))
	junitXMLFile, err := os.OpenFile(filepath.Join(outputDir, "junit_01.xml"), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("unable to create junit_01.xml: %v\n", err)
	}
	defer junitXMLFile.Close()
	err = downloadFile(config, clientset, common.PodName, common.OutputContainer, "/tmp/results/junit_01.xml", junitXMLFile)
	if err != nil {
		log.Fatalf("unable to download junit_01.xml: %v\n", err)
	}
}

// NewClient returns a new client
func NewClient() *Client {
	return &Client{}
}
