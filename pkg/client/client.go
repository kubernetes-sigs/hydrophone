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

	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/log"
)

// Client is a struct that holds the clientset and exit code
type Client struct {
	ClientSet *kubernetes.Clientset
	ExitCode  int
}

// NewClient returns a new client
func NewClient() *Client {
	return &Client{}
}

// FetchFiles downloads the e2e.log and junit_01.xml files from the pod
// and writes them to the output directory
func (c *Client) FetchFiles(ctx context.Context, config *rest.Config, clientset *kubernetes.Clientset, outputDir string) error {
	if err := c.fetchFile(ctx, config, clientset, outputDir, "e2e.log"); err != nil {
		return err
	}

	if err := c.fetchFile(ctx, config, clientset, outputDir, "junit_01.xml"); err != nil {
		return err
	}

	return nil
}

// FetchFiles downloads a single file from the output container to the local machine.
func (c *Client) fetchFile(ctx context.Context, config *rest.Config, clientset *kubernetes.Clientset, outputDir string, filename string) error {
	dest := filepath.Join(outputDir, filename)
	log.Printf("Downloading %s to %s...", filename, dest)

	localFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer localFile.Close()

	containerFile := "/tmp/results/" + filename

	return downloadFile(ctx, config, clientset, viper.GetString("namespace"), common.PodName, common.OutputContainer, containerFile, localFile)
}
