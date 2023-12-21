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
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	v1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	ctx = context.Background()
)

// Initializes the kube config clientset
func Init(cfg *ArgConfig) (*rest.Config, *kubernetes.Clientset) {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := getKubeConfig(cfg.Kubeconfig)

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("kubeconfig can't be loaded: %v\n", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("error getting config client: %v\n", err)
	}

	return config, clientset
}

func getKubeConfig(kubeconfig string) string {
	homeDir := os.Getenv("HOME")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(homeDir, ".kube", "config")
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}
	}

	// Handle cases where kubeconfig is set to users home directory in linux
	if strings.HasPrefix(kubeconfig, "~") {
		kubeconfig = filepath.Join(homeDir, kubeconfig[1:])
	}

	return kubeconfig
}

func RunE2E(clientset *kubernetes.Clientset, cfg *ArgConfig) {
	conformanceNS := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: Namespace,
		},
	}

	conformanceSA := v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name:      ServiceAccountName,
			Namespace: "conformance",
		},
	}

	conformanceClusterRole := rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name: ClusterRoleName,
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
			{
				NonResourceURLs: []string{"/metrics", "/logs", "/logs/*"},
				Verbs:           []string{"get"},
			},
		},
	}

	conformanceClusterRoleBinding := rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name: ClusterRoleBindingName,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "conformance-serviceaccount",
		},
		Subjects: []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "conformance-serviceaccount",
				Namespace: "conformance",
			},
		},
	}

	conformancePod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "e2e-conformance-test",
			Namespace: "conformance",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            ConformanceContainer,
					Image:           cfg.Image,
					ImagePullPolicy: v1.PullIfNotPresent,
					Env: []v1.EnvVar{
						{
							Name:  "E2E_FOCUS",
							Value: fmt.Sprintf("%s", cfg.Focus),
						},
						{
							Name:  "E2E_SKIP",
							Value: fmt.Sprintf("%s", cfg.Skip),
						},
						{
							Name:  "E2E_PROVIDER",
							Value: "skeleton",
						},
						{
							Name:  "E2E_PARALLEL",
							Value: fmt.Sprintf("%d", cfg.Parallel),
						},
						{
							Name:  "E2E_VERBOSITY",
							Value: fmt.Sprintf("%d", cfg.Verbosity),
						},
						{
							Name:  "E2E_USE_GO_RUNNER",
							Value: "true",
						},
					},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      "output-volume",
							MountPath: "/tmp/results",
						},
					},
				},
				{
					Name:    OutputContainer,
					Image:   "busybox",
					Command: []string{"/bin/sh", "-c", "sleep infinity"},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      "output-volume",
							MountPath: "/tmp/results",
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "output-volume",
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{},
					},
				},
			},
			RestartPolicy:      v1.RestartPolicyNever,
			ServiceAccountName: "conformance-serviceaccount",
		},
	}

	ns, err := clientset.CoreV1().Namespaces().Create(ctx, &conformanceNS, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Printf("namespace already exist %s", PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("namespace created %s\n", ns.Name)

	sa, err := clientset.CoreV1().ServiceAccounts(ns.Name).Create(ctx, &conformanceSA, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Printf("serviceaccount already exist %s", PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("serviceaccount created %s\n", sa.Name)

	clusterRole, err := clientset.RbacV1().ClusterRoles().Create(ctx, &conformanceClusterRole, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Printf("clusterrole already exist %s", PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("clusterrole created %s\n", clusterRole.Name)

	clusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Create(ctx, &conformanceClusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Printf("clusterrolebinding already exist %s", PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("clusterrolebinding created %s\n", clusterRoleBinding.Name)

	pod, err := clientset.CoreV1().Pods(ns.Name).Create(ctx, &conformancePod, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Printf("pod already exist %s", PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("pod created %s\n", pod.Name)
}

func Cleanup(clientset *kubernetes.Clientset) {
	err := clientset.CoreV1().Pods(Namespace).Delete(ctx, PodName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("pod %s doesn't exist\n", PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("pod deleted %s\n", PodName)

	err = clientset.RbacV1().ClusterRoleBindings().Delete(ctx, ClusterRoleBindingName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("clusterrolebinding %s doesn't exist\n", ClusterRoleBindingName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("clusterrolebinding deleted %s\n", ClusterRoleBindingName)

	err = clientset.RbacV1().ClusterRoles().Delete(ctx, ClusterRoleName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("clusterrole %s doesn't exist\n", ClusterRoleName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("clusterrole deleted %s\n", ClusterRoleName)

	err = clientset.CoreV1().ServiceAccounts(Namespace).Delete(ctx, ServiceAccountName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("serviceaccount %s doesn't exist\n", ServiceAccountName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("serviceaccount deleted %s\n", ServiceAccountName)

	err = clientset.CoreV1().Namespaces().Delete(ctx, Namespace, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("namespace %s doesn't exist\n", Namespace)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("namespace deleted %s\n", Namespace)
}
