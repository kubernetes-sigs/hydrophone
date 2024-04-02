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
	"bytes"
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func downloadFile(ctx context.Context, config *rest.Config, clientset *kubernetes.Clientset,
	namespace, podName, containerName, filePath string,
	writer io.Writer) error {
	// Create an exec request
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", containerName)

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return err
	}
	// Configure exec options
	option := &corev1.PodExecOptions{
		Stdout:  true,
		Stderr:  true,
		Command: []string{"cat", filePath},
	}
	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(option, parameterCodec)

	// Create an executor
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	// Stream the file content from the container to the writer
	var stderr bytes.Buffer

	err = exec.StreamWithContext(
		ctx,
		remotecommand.StreamOptions{
			Stdout: writer,
			Stderr: &stderr,
		})
	if err != nil {
		return fmt.Errorf("download failed: %w (stderr: %s)", err, stderr.String())
	}

	return nil
}
