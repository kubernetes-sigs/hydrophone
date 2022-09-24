package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	container_image = "registry.k8s.io/conformance-amd64:v1.25.0"
)

var (
	ctx = context.Background()
)

// Initializes the kube confif clientset
func Init() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		homeDir := os.Getenv("HOME")
		kubeconfig := filepath.Join(homeDir, ".kube", "config")
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("kubeconfig can't be loaded: %v\n", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("error getting config client: %v\n", err)
	}

	return clientset
}

func RunE2E(clientset *kubernetes.Clientset, focus string) {
	conformanceNS := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "conformance",
		},
	}

	conformanceSA := v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name:      "conformance-serviceaccount",
			Namespace: "conformance",
		},
	}

	conformanceClusterRole := rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"component": "conformance",
			},
			Name: "conformance-serviceaccount",
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
			Name: "conformance-serviceaccount-role",
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
					Name:            "conformance-container",
					Image:           container_image,
					ImagePullPolicy: v1.PullIfNotPresent,
					Env: []v1.EnvVar{
						{
							Name:  "E2E_FOCUS",
							Value: fmt.Sprintf("\\[%s\\]", focus),
						},
						{
							Name:  "E2E_SKIP",
							Value: "",
						},
						{
							Name:  "E2E_PROVIDER",
							Value: "skeleton",
						},
						{
							Name:  "E2E_PARALLEL",
							Value: "false",
						},
						{
							Name:  "E2E_VERBOSITY",
							Value: "4",
						},
					},
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
						HostPath: &v1.HostPathVolumeSource{
							Path: "/tmp/results",
						},
					},
				},
			},
			RestartPolicy:      v1.RestartPolicyNever,
			ServiceAccountName: "conformance-serviceaccount",
		},
	}

	ns, err := clientset.CoreV1().Namespaces().Create(ctx, &conformanceNS, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("namespace created %s\n", ns.Name)

	sa, err := clientset.CoreV1().ServiceAccounts(ns.Name).Create(ctx, &conformanceSA, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("serviceaccount created %s\n", sa.Name)

	clusterRole, err := clientset.RbacV1().ClusterRoles().Create(ctx, &conformanceClusterRole, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("clusterrole created %s\n", clusterRole.Name)

	clusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Create(ctx, &conformanceClusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("clusterrolebinding created %s\n", clusterRoleBinding.Name)

	pod, err := clientset.CoreV1().Pods(ns.Name).Create(ctx, &conformancePod, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("pod created %s\n", pod.Name)
}
