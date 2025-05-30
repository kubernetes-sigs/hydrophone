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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client is used to retrieve conformance test logs, results and status information.
type Client struct {
	config    *rest.Config
	clientset *kubernetes.Clientset
	namespace string
}

// NewClient creates a client for interacting with the conformance test pod
func NewClient(config *rest.Config, clientset *kubernetes.Clientset, namespace string) *Client {
	return &Client{
		config:    config,
		clientset: clientset,
		namespace: namespace,
	}
}
