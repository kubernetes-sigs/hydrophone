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
	"slices"
	"time"

	"sigs.k8s.io/hydrophone/pkg/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreatePod creates a new Pod and waits for it to be Running/Succeeded.
func CreatePod(ctx context.Context, cs *kubernetes.Clientset, pod *corev1.Pod, timeout time.Duration) (*corev1.Pod, error) {
	// Create the pod in the cluster
	created, err := cs.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create Pod: %w", err)
	}

	log.Printf("Created Pod %s.", pod.Name)

	// From now on, even in errors return the created Pod
	// to allow the caller to perform cleanups if desired.

	// Watch for pod events
	watcher, err := cs.CoreV1().Pods(created.Namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + created.Name,
	})
	if err != nil {
		return created, fmt.Errorf("failed to watch Pod events: %w", err)
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

			if err := CheckFailedPod(pod); err != nil {
				return created, err
			}

			if pod.Status.Phase != corev1.PodPending {
				return created, nil
			}

			lastStatus = &pod.Status

		case <-deadline:
			phase := corev1.PodUnknown
			if lastStatus != nil {
				phase = lastStatus.Phase
			}

			return created, fmt.Errorf("timed out waiting for Pod, last status was %v", phase)
		}
	}
}

// containerErrorReasons is a list of possible reasons a container in a Pod can have.
// If a container has this status, CreatePod() considers it fails and aborts.
var containerErrorReasons = []string{"ErrImagePull", "ImagePullBackOff", "Error", "CrashLoopBackOff"}

func CheckFailedPod(pod *corev1.Pod) error {
	for _, cs := range pod.Status.ContainerStatuses {
		if s := cs.State.Waiting; s != nil && slices.Contains(containerErrorReasons, s.Reason) {
			return errors.New(s.Message)
		}

		if s := cs.State.Terminated; s != nil && slices.Contains(containerErrorReasons, s.Reason) {
			msg := s.Message
			if msg == "" {
				msg = "Pod has encountered an error"
			}

			return errors.New(msg)
		}
	}

	return nil
}
