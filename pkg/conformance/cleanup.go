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

package conformance

import (
	"context"
	"fmt"
	"time"

	"sigs.k8s.io/hydrophone/pkg/log"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// Cleanup removes all resources created during E2E tests.
func (r *TestRunner) Cleanup(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	name := r.namespacedName(ClusterRoleBindingName)
	err := r.clientset.RbacV1().ClusterRoleBindings().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		log.Printf("Deleted ClusterRoleBinding %s.", name)
	}

	name = r.namespacedName(ClusterRoleName)
	err = r.clientset.RbacV1().ClusterRoles().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		log.Printf("Deleted ClusterRole %s.", name)
	}

	name = r.config.Namespace
	err = r.clientset.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}

	log.Printf("Waiting for Namespace %s to be deleted (timeout %v).", name, r.config.CleanupTimeout)

	// Poll instead of watching: watch streams get closed by the API server and
	// restarting them can miss the deletion event. Polling cannot.
	err = wait.PollUntilContextTimeout(ctx, 10*time.Second, r.config.CleanupTimeout, true, func(ctx context.Context) (bool, error) {
		_, err := r.clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			return true, nil
		}

		if err != nil && ctx.Err() == nil {
			log.Printf("Error checking Namespace %s: %v (retrying)", name, err)
		}

		return false, nil
	})
	if err != nil {
		return fmt.Errorf("namespace %s still exists after %v: %w", name, r.config.CleanupTimeout, err)
	}

	log.Printf("Deleted Namespace %s.", name)

	return nil
}
