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
	"github.com/dims/hydrophone/pkg/service"
	"k8s.io/client-go/rest"
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
)

var (
	ctx = context.TODO()
)

type Client struct {
	ClientSet *kubernetes.Clientset
	ExitCode  int
}

func (c *Client) FetchFiles(config *rest.Config, clientset *kubernetes.Clientset) {
	log.Println("downloading e2e.log")
	e2eLogFile, err := os.OpenFile("e2e.log", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("unable to create e2e.log: %v\n", err)
	}
	defer e2eLogFile.Close()
	err = downloadFile(config, clientset, service.PodName, service.OutputContainer, "/tmp/results/e2e.log", e2eLogFile)
	if err != nil {
		log.Fatalf("unable to download e2e.log: %v\n", err)
	}
	log.Println("downloading junit_01.xml")
	junitXMLFile, err := os.OpenFile("junit_01.xml", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("unable to create junit_01.xml: %v\n", err)
	}
	defer junitXMLFile.Close()
	err = downloadFile(config, clientset, service.PodName, service.OutputContainer, "/tmp/results/junit_01.xml", junitXMLFile)
	if err != nil {
		log.Fatalf("unable to download junit_01.xml: %v\n", err)
	}
}

// Return a new Client
func NewClient() *Client {
	return &Client{}
}
