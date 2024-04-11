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
	"errors"

	"sigs.k8s.io/hydrophone/pkg/log"

	v1 "k8s.io/api/core/v1"
)

// Exit process if conformance pod happen ImagePullBackOff.
func ExitWhenImagePullBackOff(pod *v1.Pod) {
	if pod.Status.Phase != v1.PodPending {
		return
	}
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.State.Waiting != nil && cs.State.Waiting.Reason == "ImagePullBackOff" {
			log.Fatal(errors.New(cs.State.Waiting.Message))
		}
	}
}
