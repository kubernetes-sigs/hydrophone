/*
Copyright 2024 The Kubernetes Authors.

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
	"context"
	"errors"
	"fmt"
	"time"

	"sigs.k8s.io/hydrophone/pkg/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreatePod creates a new Pod and waits for it to be Running.
func CreatePod(ctx context.Context, cs *kubernetes.Clientset, pod *corev1.Pod, timeout time.Duration) (*corev1.Pod, error) {
	// Create the pod in the cluster
	created, err := cs.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create Pod: %w", err)
	}

	// Watch for pod events
	watcher, err := cs.CoreV1().Pods(created.Namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + created.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to watch Pod events: %w", err)
	}
	defer watcher.Stop()

	log.Printf("Waiting up to %v for Pod to start...", timeout)

	deadline := time.After(timeout)
	var lastStatus *corev1.PodStatus

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return created, nil
			}

			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue
			}

			lastStatus = &pod.Status

			if lastStatus.Phase == corev1.PodRunning {
				return created, nil
			}

			for _, cs := range lastStatus.ContainerStatuses {
				if cs.State.Waiting != nil && cs.State.Waiting.Reason == "ImagePullBackOff" {
					return nil, errors.New(cs.State.Waiting.Message)
				}
			}

		case <-deadline:
			phase := corev1.PodUnknown
			if lastStatus != nil {
				phase = lastStatus.Phase
			}

			return nil, fmt.Errorf("timed out waiting for Pod, last status was %v", phase)
		}
	}
}
