package service

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
