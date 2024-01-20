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
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/log"
)

var (
	ctx = context.Background()
)

// Init Initializes the kube config clientset
func Init(kubeconfig string) (*rest.Config, *kubernetes.Clientset) {
	config, err := rest.InClusterConfig()
	if err != nil {
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

func GetKubeConfig(kubeconfig string) string {
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

// RunE2E sets up the necessary resources and runs E2E conformance tests.
func RunE2E(clientset *kubernetes.Clientset) {
	conformanceNS := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: common.Namespace,
		},
	}

	conformanceSA := v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name:      common.ServiceAccountName,
			Namespace: "conformance",
		},
	}

	conformanceClusterRole := rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name: common.ClusterRoleName,
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
			Name: common.ClusterRoleBindingName,
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
					Name:            common.ConformanceContainer,
					Image:           viper.GetString("conformance-image"),
					ImagePullPolicy: v1.PullIfNotPresent,
					Env: []v1.EnvVar{
						{
							Name:  "E2E_FOCUS",
							Value: fmt.Sprintf("%s", viper.Get("focus")),
						},
						{
							Name:  "E2E_SKIP",
							Value: fmt.Sprintf("%s", viper.Get("skip")),
						},
						{
							Name:  "E2E_PROVIDER",
							Value: "skeleton",
						},
						{
							Name:  "E2E_PARALLEL",
							Value: fmt.Sprintf("%d", viper.Get("parallel")),
						},
						{
							Name:  "E2E_VERBOSITY",
							Value: fmt.Sprintf("%d", viper.Get("verbosity")),
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
					Name:    common.OutputContainer,
					Image:   viper.GetString("busybox-image"),
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

	if viper.GetBool("dry-run") {
		conformancePod.Spec.Containers[0].Env = append(conformancePod.Spec.Containers[0].Env, DryRun())
	}

	ns, err := clientset.CoreV1().Namespaces().Create(ctx, &conformanceNS, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Printf("namespace already exist %s", common.PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("namespace created %s\n", ns.Name)

	sa, err := clientset.CoreV1().ServiceAccounts(ns.Name).Create(ctx, &conformanceSA, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Printf("serviceaccount already exist %s", common.PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("serviceaccount created %s\n", sa.Name)

	clusterRole, err := clientset.RbacV1().ClusterRoles().Create(ctx, &conformanceClusterRole, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Printf("clusterrole already exist %s", common.PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("clusterrole created %s\n", clusterRole.Name)

	clusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Create(ctx, &conformanceClusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Printf("clusterrolebinding already exist %s", common.PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("clusterrolebinding created %s\n", clusterRoleBinding.Name)

	if viper.GetString("test-repo-list") != "" {
		RepoListData, err := os.ReadFile(viper.GetString("test-repo-list"))
		if err != nil {
			log.Fatal(err)
		}
		configMap := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "repo-list-config",
				Namespace: common.Namespace,
			},
			Data: map[string]string{
				"repo-list.yaml": string(RepoListData),
			},
		}

		conformancePod.Spec.Volumes = append(conformancePod.Spec.Volumes,
			v1.Volume{
				Name: "repo-list-volume",
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "repo-list-config",
						},
					},
				},
			})

		conformancePod.Spec.Containers[0].VolumeMounts = append(conformancePod.Spec.Containers[0].VolumeMounts,
			v1.VolumeMount{
				Name:      "repo-list-volume",
				MountPath: "/tmp/repo-list",
				ReadOnly:  true,
			})

		conformancePod.Spec.Containers[0].Env = append(conformancePod.Spec.Containers[0].Env, v1.EnvVar{
			Name:  "KUBE_TEST_REPO_LIST",
			Value: "/tmp/repo-list/repo-list.yaml",
		})

		cm, err := clientset.CoreV1().ConfigMaps(common.Namespace).Create(ctx, configMap, metav1.CreateOptions{})
		if err != nil {
			if errors.IsAlreadyExists(err) {
				log.Printf("configmap already exists %s", configMap.ObjectMeta.Name)
			} else {
				log.Fatal(err)
			}
		}
		log.Printf("configmap created %s\n", cm.Name)
	}

	if viper.GetString("test-repo") != "" {
		conformancePod.Spec.Containers[0].Env = append(conformancePod.Spec.Containers[0].Env, v1.EnvVar{
			Name:  "KUBE_TEST_REPO",
			Value: viper.GetString("test-repo"),
		})
	}

	pod, err := clientset.CoreV1().Pods(ns.Name).Create(ctx, &conformancePod, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Printf("pod already exist %s", common.PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("pod created %s\n", pod.Name)
}

// Cleanup removes all resources created during E2E tests.
func Cleanup(clientset *kubernetes.Clientset) {
	err := clientset.CoreV1().Pods(common.Namespace).Delete(ctx, common.PodName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("pod %s doesn't exist\n", common.PodName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("pod deleted %s\n", common.PodName)

	err = clientset.RbacV1().ClusterRoleBindings().Delete(ctx, common.ClusterRoleBindingName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("clusterrolebinding %s doesn't exist\n", common.ClusterRoleBindingName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("clusterrolebinding deleted %s\n", common.ClusterRoleBindingName)

	err = clientset.RbacV1().ClusterRoles().Delete(ctx, common.ClusterRoleName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("clusterrole %s doesn't exist\n", common.ClusterRoleName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("clusterrole deleted %s\n", common.ClusterRoleName)

	err = clientset.CoreV1().ServiceAccounts(common.Namespace).Delete(ctx, common.ServiceAccountName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("serviceaccount %s doesn't exist\n", common.ServiceAccountName)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("serviceaccount deleted %s\n", common.ServiceAccountName)

	err = clientset.CoreV1().Namespaces().Delete(ctx, common.Namespace, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("namespace %s doesn't exist\n", common.Namespace)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("namespace deleted %s\n", common.Namespace)
}

// DryRun returns an environment variable to tell the conformance test to run in dry run mode.
func DryRun() v1.EnvVar {
	return v1.EnvVar{
		Name:  "E2E_DRYRUN",
		Value: "true",
	}
}
