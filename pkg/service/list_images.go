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

package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/log"

	"github.com/spf13/viper"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PrintListImages creates and runs a conformance image with the --list-images flag
// This will print a list of all the images used by the conformance image.
func PrintListImages(ctx context.Context, clientSet *kubernetes.Clientset) error {
	// Create a pod object definition
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "list-images-",
			Namespace:    "default",
			Annotations: map[string]string{
				"list-images": "true",
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyOnFailure,
			Containers: []corev1.Container{
				{
					Name:  common.ConformanceContainer,
					Image: viper.GetString("conformance-image"),
					Command: []string{
						"/usr/local/bin/e2e.test",
						"--list-images",
					},
				},
			},
		},
	}

	// Create the pod in the cluster
	pod, err := clientSet.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Pod: %w", err)
	}
	log.Printf("Created Pod %s.", pod.Name)

	// Watch for pod events
	watcher, err := clientSet.CoreV1().Pods("default").Watch(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + pod.Name,
	})
	if err != nil {
		return fmt.Errorf("failed to watch Pod events: %w", err)
	}
	defer watcher.Stop()

	log.Println("Waiting for Pod to complete...")

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return nil
			}

			// Handle pod event
			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue
			}

			// Check if the pod is in a terminal state
			if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
				// Trigger desired action (e.g., fetching and printing logs)
				log.Printf("Pod completed: %s", pod.Status.Phase)

				// Fetch the logs
				req := clientSet.CoreV1().Pods("default").GetLogs(pod.Name, &corev1.PodLogOptions{})
				podLogs, err := req.Stream(ctx)
				if err != nil {
					return fmt.Errorf("failed to fetch Pod logs: %w", err)
				}
				defer podLogs.Close()

				// Read and print the logs
				buf := new(bytes.Buffer)
				_, err = io.Copy(buf, podLogs)
				if err != nil {
					return fmt.Errorf("failed to read Pod logs: %w", err)
				}

				lines := strings.Split(buf.String(), "\n")
				sort.Strings(lines)
				for _, line := range lines {
					fmt.Println(line)
				}

				err = clientSet.CoreV1().Pods("default").Delete(ctx, pod.Name, metav1.DeleteOptions{})
				if err != nil {
					return fmt.Errorf("failed to delete Pod: %w", err)
				}
				return nil
			}

		case <-time.After(2 * time.Second):
			// Check status every 2 seconds
		}
	}
}
