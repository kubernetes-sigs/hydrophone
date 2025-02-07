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
	"os"
	"strings"
	"time"

	"k8s.io/utils/ptr"

	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/log"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Deploy sets up the necessary resources and runs E2E conformance tests.
func (r *TestRunner) Deploy(ctx context.Context, focus, skipPreflight string, verboseGinkgo bool, timeout time.Duration) error {
	conformanceNS := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.config.Namespace,
		},
	}

	conformanceSA := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name:      ServiceAccountName,
			Namespace: r.config.Namespace,
		},
	}

	conformanceClusterRole := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name: r.namespacedName(ClusterRoleName),
		},
		Rules: []rbacv1.PolicyRule{
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

	conformanceClusterRoleBinding := rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name: r.namespacedName(ClusterRoleBindingName),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     r.namespacedName(ClusterRoleName),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      ServiceAccountName,
				Namespace: r.config.Namespace,
			},
		},
	}

	containerEnv := []corev1.EnvVar{
		{
			Name:  "E2E_FOCUS",
			Value: focus,
		},
		{
			Name:  "E2E_SKIP",
			Value: r.config.Skip,
		},
		{
			Name:  "E2E_PROVIDER",
			Value: "skeleton",
		},
		{
			Name:  "E2E_VERBOSITY",
			Value: fmt.Sprintf("%d", r.config.Verbosity),
		},
		{
			Name:  "E2E_USE_GO_RUNNER",
			Value: "true",
		},
		{
			Name:  "E2E_EXTRA_ARGS",
			Value: strings.Join(r.config.ExtraArgs, " "),
		},
	}

	extraGinkgoArgs := r.config.ExtraGinkgoArgs

	if r.config.Parallel > 1 {
		extraGinkgoArgs = append(extraGinkgoArgs, fmt.Sprintf("--procs=%d", r.config.Parallel))
	}

	if verboseGinkgo {
		extraGinkgoArgs = append(extraGinkgoArgs, "-v")
	}

	containerEnv = append(containerEnv, corev1.EnvVar{
		Name:  "E2E_EXTRA_GINKGO_ARGS",
		Value: strings.Join(extraGinkgoArgs, " "),
	})

	if r.config.DryRun {
		containerEnv = append(containerEnv, corev1.EnvVar{
			Name:  "E2E_DRYRUN",
			Value: "true",
		})
	}

	conformancePod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "e2e-conformance-test",
			Namespace: conformanceNS.Name,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            ConformanceContainer,
					Image:           r.config.ConformanceImage,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Env:             containerEnv,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "output-volume",
							MountPath: "/tmp/results",
						},
					},
					SecurityContext: &corev1.SecurityContext{
						AllowPrivilegeEscalation: ptr.To(false),
						Capabilities: &corev1.Capabilities{
							Drop: []corev1.Capability{
								"ALL",
							},
						},
						RunAsNonRoot: ptr.To(true),
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
						RunAsUser: ptr.To(int64(65534)),
					},
				},
				{
					Name:    OutputContainer,
					Image:   r.config.BusyboxImage,
					Command: []string{"/bin/sh", "-c", "sleep infinity"},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "output-volume",
							MountPath: "/tmp/results",
						},
					},
					SecurityContext: &corev1.SecurityContext{
						AllowPrivilegeEscalation: ptr.To(false),
						Capabilities: &corev1.Capabilities{
							Drop: []corev1.Capability{
								"ALL",
							},
						},
						RunAsNonRoot: ptr.To(true),
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
						RunAsUser: ptr.To(int64(65534)),
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "output-volume",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
			RestartPolicy:      corev1.RestartPolicyNever,
			ServiceAccountName: ServiceAccountName,
			Tolerations: []corev1.Toleration{
				{
					// An empty key with operator Exists matches all keys,
					// values and effects which means this will tolerate everything.
					// As noted in https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/
					Operator: "Exists",
				},
			},
		},
	}

	ns, err := r.clientset.CoreV1().Namespaces().Create(ctx, &conformanceNS, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			if skipPreflight != "" {
				log.Printf("Using existing namespace: %s", r.config.Namespace)
			} else {
				//nolint:stylecheck // error message references a Kubernetes resource type.
				return fmt.Errorf("namespace %s already exists, please run with --cleanup first", conformanceNS.Name)
			}
		} else {
			return fmt.Errorf("failed to create namespace: %w", err)
		}
	} else {
		log.Printf("Created namespace %s.", ns.Name)
	}

	sa, err := r.clientset.CoreV1().ServiceAccounts(r.config.Namespace).Create(ctx, &conformanceSA, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			if skipPreflight != "" {
				log.Printf("using existing ServiceAccount: %s/%s", r.config.Namespace, ServiceAccountName)
			} else {
				//nolint:stylecheck // error message references a Kubernetes resource type.
				return fmt.Errorf("ServiceAccount %s already exist, please run --cleanup first", conformanceSA.Name)
			}
		} else {
			return fmt.Errorf("failed to create ServiceAccount: %w", err)
		}
	} else {
		log.Printf("Created ServiceAccount %s.", sa.Name)
	}

	clusterRole, err := r.clientset.RbacV1().ClusterRoles().Create(ctx, &conformanceClusterRole, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			if skipPreflight != "" {
				log.Printf("using existing ClusterRole: %s/%s", r.config.Namespace, r.namespacedName(ClusterRoleName))
			} else {
				//nolint:stylecheck // error message references a Kubernetes resource type.
				return fmt.Errorf("ClusterRole %s already exist, please run --cleanup first", conformanceClusterRole.Name)
			}
		} else {
			return fmt.Errorf("failed to create ClusterRole: %w", err)
		}
	} else {
		log.Printf("Created ClusterRole %s.", clusterRole.Name)
	}

	clusterRoleBinding, err := r.clientset.RbacV1().ClusterRoleBindings().Create(ctx, &conformanceClusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			if skipPreflight != "" {
				log.Printf("using existing ClusterRoleBinding: %s/%s", r.config.Namespace, r.namespacedName(ClusterRoleBindingName))
			} else {
				//nolint:stylecheck // error message references a Kubernetes resource type.
				return fmt.Errorf("ClusterRoleBinding %s already exist, please run --cleanup first", conformanceClusterRoleBinding.Name)
			}
		} else {
			return fmt.Errorf("failed to create ClusterRoleBinding: %w", err)
		}
	} else {
		log.Printf("Created ClusterRoleBinding %s.", clusterRoleBinding.Name)
	}

	if filename := r.config.TestRepoList; filename != "" {
		repoListData, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read repo list: %w", err)
		}

		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "repo-list-config",
				Namespace: r.config.Namespace,
			},
			Data: map[string]string{
				"repo-list.yaml": string(repoListData),
			},
		}

		conformancePod.Spec.Volumes = append(conformancePod.Spec.Volumes,
			corev1.Volume{
				Name: "repo-list-volume",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "repo-list-config",
						},
					},
				},
			})

		conformancePod.Spec.Containers[0].VolumeMounts = append(conformancePod.Spec.Containers[0].VolumeMounts,
			corev1.VolumeMount{
				Name:      "repo-list-volume",
				MountPath: "/tmp/repo-list",
				ReadOnly:  true,
			})

		conformancePod.Spec.Containers[0].Env = append(conformancePod.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "KUBE_TEST_REPO_LIST",
			Value: "/tmp/repo-list/repo-list.yaml",
		})

		cm, err := r.clientset.CoreV1().ConfigMaps(ns.Name).Create(ctx, configMap, metav1.CreateOptions{})
		if err != nil {
			if errors.IsAlreadyExists(err) {
				//nolint:stylecheck // error message references a Kubernetes resource type.
				err = fmt.Errorf("ConfigMap %s already exist, please run --cleanup first", configMap.Name)
			}

			return err
		}
		log.Printf("Created ConfigMap %s.", cm.Name)
	}

	if r.config.TestRepo != "" {
		conformancePod.Spec.Containers[0].Env = append(conformancePod.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "KUBE_TEST_REPO",
			Value: r.config.TestRepo,
		})
	}

	pod, err := common.CreatePod(ctx, r.clientset, &conformancePod, timeout)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			if skipPreflight != "" {
				log.Printf("using existing Pod: %s/%s", r.config.Namespace, "e2e-conformance-test")
			} else {
				//nolint:stylecheck // error message references a Kubernetes resource type.
				return fmt.Errorf("Pod %s already exist, please run --cleanup first", conformanceClusterRoleBinding.Name)
			}
		} else {
			return fmt.Errorf("failed to create Pod: %w", err)
		}
	} else {
		log.Printf("Created ConformancePod %s.", pod.Name)
	}

	return nil
}
