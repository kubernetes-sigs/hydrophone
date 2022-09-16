package client

import (
	"bufio"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// List pod resource with the given namespace
func GetPodLogs(clientset *kubernetes.Clientset) error {
	podLogOpts := v1.PodLogOptions{}

	req := clientset.CoreV1().Pods(NAMESPACE).GetLogs(POD_NAME, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return err
	}
	defer podLogs.Close()

	reader := bufio.NewScanner(podLogs)
	for {
		for reader.Scan() {
			line := reader.Text()
			fmt.Println(line)
		}
	}
}
