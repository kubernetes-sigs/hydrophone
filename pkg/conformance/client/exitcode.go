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
	"fmt"

	"sigs.k8s.io/hydrophone/pkg/conformance"
	"sigs.k8s.io/hydrophone/pkg/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FetchExitCode waits for pod to be in terminated state and get the exit code
func (c *Client) FetchExitCode(ctx context.Context) (int, error) {
	// Watching the pod's status
	watchInterface, err := c.clientset.CoreV1().Pods(c.namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", conformance.PodName),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to watch Pods: %w", err)
	}

	log.Println("Waiting for Pod to terminate...")
	for event := range watchInterface.ResultChan() {
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			log.Printf("Received unexpected %T object from Watch.", pod)
			return -1, nil
		}

		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			log.Println("Pod terminated.")
			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.Name == conformance.ConformanceContainer && containerStatus.State.Terminated != nil {
					return int(containerStatus.State.Terminated.ExitCode), nil
				}
			}

			return -1, fmt.Errorf("%s Pod is not terminated.", conformance.ConformanceContainer)
		}

		if pod.Status.Phase == corev1.PodRunning {
			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.State.Terminated != nil {
					log.Printf("Container %s terminated.", containerStatus.Name)
					if containerStatus.Name == conformance.ConformanceContainer {
						return int(containerStatus.State.Terminated.ExitCode), nil
					}
				}
			}
		}
	}

	return -1, nil
}
