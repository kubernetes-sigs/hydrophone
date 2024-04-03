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

// Init Initializes the kube config clientset
func Init(kubeconfig string) (*rest.Config, *kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, nil, fmt.Errorf("error loading kubeconfig: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting config client: %w", err)
	}

	return config, clientset, nil
}

// GetKubeConfig returns the path to the Kubernetes configuration file
func GetKubeConfig(kubeconfig string) string {
	homeDir := os.Getenv("HOME")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(homeDir, ".kube", "config")
		if envvar := os.Getenv("KUBECONFIG"); envvar != "" {
			kubeconfig = envvar
		}
	}

	// Handle cases where kubeconfig is set to users home directory in linux
	if strings.HasPrefix(kubeconfig, "~") {
		kubeconfig = filepath.Join(homeDir, kubeconfig[1:])
	}

	return kubeconfig
}

func namespacedName(basename string) string {
	return fmt.Sprintf("%s:%s", basename, viper.GetString("namespace"))
}

// RunE2E sets up the necessary resources and runs E2E conformance tests.
func RunE2E(ctx context.Context, clientset *kubernetes.Clientset, verboseGinkgo bool) error {
	namespace := viper.GetString("namespace")

	conformanceNS := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	conformanceSA := v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name:      common.ServiceAccountName,
			Namespace: namespace,
		},
	}

	conformanceClusterRole := rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name: namespacedName(common.ClusterRoleName),
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
			Name: namespacedName(common.ClusterRoleBindingName),
		},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     namespacedName(common.ClusterRoleName),
		},
		Subjects: []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      common.ServiceAccountName,
				Namespace: namespace,
			},
		},
	}

	containerEnv := []v1.EnvVar{
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
			Name:  "E2E_VERBOSITY",
			Value: fmt.Sprintf("%d", viper.Get("verbosity")),
		},
		{
			Name:  "E2E_USE_GO_RUNNER",
			Value: "true",
		},
		{
			Name:  "E2E_EXTRA_ARGS",
			Value: strings.Join(viper.GetStringSlice("extra-args"), " "),
		},
	}

	extraGinkgoArgs := viper.GetStringSlice("extra-ginkgo-args")

	if threads := viper.GetInt("parallel"); threads > 1 {
		extraGinkgoArgs = append(extraGinkgoArgs, fmt.Sprintf("--procs=%d", threads))
		containerEnv = append(containerEnv, v1.EnvVar{
			Name:  "E2E_PARALLEL",
			Value: "true",
		})
	}

	if verboseGinkgo {
		extraGinkgoArgs = append(extraGinkgoArgs, "-v")
	}

	containerEnv = append(containerEnv, v1.EnvVar{
		Name:  "E2E_EXTRA_GINKGO_ARGS",
		Value: strings.Join(extraGinkgoArgs, " "),
	})

	if viper.GetBool("dry-run") {
		containerEnv = append(containerEnv, v1.EnvVar{
			Name:  "E2E_DRYRUN",
			Value: "true",
		})
	}

	conformancePod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "e2e-conformance-test",
			Namespace: conformanceNS.Name,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            common.ConformanceContainer,
					Image:           viper.GetString("conformance-image"),
					ImagePullPolicy: v1.PullIfNotPresent,
					Env:             containerEnv,
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
			ServiceAccountName: common.ServiceAccountName,
			Tolerations: []v1.Toleration{
				{
					// An empty key with operator Exists matches all keys,
					// values and effects which means this will tolerate everything.
					// As noted in https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/
					Operator: "Exists",
				},
			},
		},
	}

	ns, err := clientset.CoreV1().Namespaces().Create(ctx, &conformanceNS, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			//nolint:stylecheck // error message references a Kubernetes resource type.
			err = fmt.Errorf("Namespace %s already exist, please run --cleanup first", conformanceNS.Name)
		}

		return err
	}
	log.Printf("Created Namespace %s.", ns.Name)

	sa, err := clientset.CoreV1().ServiceAccounts(ns.Name).Create(ctx, &conformanceSA, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			//nolint:stylecheck // error message references a Kubernetes resource type.
			err = fmt.Errorf("ServiceAccount %s already exist, please run --cleanup first", conformanceSA.Name)
		}

		return err
	}
	log.Printf("Created ServiceAccount %s.", sa.Name)

	clusterRole, err := clientset.RbacV1().ClusterRoles().Create(ctx, &conformanceClusterRole, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			//nolint:stylecheck // error message references a Kubernetes resource type.
			err = fmt.Errorf("ClusterRole %s already exist, please run --cleanup first", conformanceClusterRole.Name)
		}

		return err
	}
	log.Printf("Created Clusterrole %s.", clusterRole.Name)

	clusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Create(ctx, &conformanceClusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			//nolint:stylecheck // error message references a Kubernetes resource type.
			err = fmt.Errorf("ClusterRoleBinding %s already exist, please run --cleanup first", conformanceClusterRoleBinding.Name)
		}

		return err
	}
	log.Printf("Created ClusterRoleBinding %s.", clusterRoleBinding.Name)

	if viper.GetString("test-repo-list") != "" {
		RepoListData, err := os.ReadFile(viper.GetString("test-repo-list"))
		if err != nil {
			return fmt.Errorf("failed to read repo list: %w", err)
		}

		configMap := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "repo-list-config",
				Namespace: namespace,
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

		cm, err := clientset.CoreV1().ConfigMaps(ns.Name).Create(ctx, configMap, metav1.CreateOptions{})
		if err != nil {
			if errors.IsAlreadyExists(err) {
				//nolint:stylecheck // error message references a Kubernetes resource type.
				err = fmt.Errorf("ConfigMap %s already exist, please run --cleanup first", configMap.Name)
			}

			return err
		}
		log.Printf("Created ConfigMap %s.", cm.Name)
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
			//nolint:stylecheck // error message references a Kubernetes resource type.
			err = fmt.Errorf("Pod %s already exist, please run --cleanup first", conformancePod.Name)
		}

		return err
	}
	log.Printf("Created Pod %s.", pod.Name)

	return nil
}

// Cleanup removes all resources created during E2E tests.
func Cleanup(ctx context.Context, clientset *kubernetes.Clientset) error {
	name := namespacedName(common.ClusterRoleBindingName)
	err := clientset.RbacV1().ClusterRoleBindings().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		log.Printf("Deleted ClusterRoleBinding %s.", name)
	}

	name = namespacedName(common.ClusterRoleName)
	err = clientset.RbacV1().ClusterRoles().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		log.Printf("Deleted ClusterRole %s.", name)
	}

	name = viper.GetString("namespace")
	err = clientset.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		log.Printf("Deleted Namespace %s.", name)
	}

	return nil
}
