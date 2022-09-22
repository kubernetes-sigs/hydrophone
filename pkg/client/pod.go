package client

import (
	"bufio"
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// List pod resource with the given namespace
func getPodLogs(cancelCtx context.Context, clientset *kubernetes.Clientset) error {
	podLogOpts := v1.PodLogOptions{
		Follow: true,
	}

	req := clientset.CoreV1().Pods(NAMESPACE).GetLogs(POD_NAME, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return err
	}
	defer podLogs.Close()

	reader := bufio.NewScanner(podLogs)

	for {
		select {
		case <-cancelCtx.Done():
			return nil
		default:
			for reader.Scan() {
				line := reader.Text()
				fmt.Println(line)
			}
		}
	}
}
