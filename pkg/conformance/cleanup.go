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

	"sigs.k8s.io/hydrophone/pkg/log"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
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

	// start a watcher before deleting the namespace
	watcher, err := r.clientset.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("metadata.name", r.config.Namespace).String(),
	})
	if err != nil {
		return err
	}

	defer watcher.Stop()

	name = r.config.Namespace
	err = r.clientset.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}

	log.Printf("Waiting for Namespace %s to be deleted.", name)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-watcher.ResultChan():
			if event.Type == watch.Error {
				return fmt.Errorf("error watching waiting for namespace deletion")
			}

			if event.Type == watch.Deleted {
				log.Printf("Deleted Namespace %s.", name)

				return nil
			}
		}
	}
}
