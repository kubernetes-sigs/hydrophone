package client

import (
	"bufio"

	"github.com/dims/k8s-run-e2e/pkg/service"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// List pod resource with the given namespace
func getPodLogs(clientset *kubernetes.Clientset, stream streamLogs) {
	podLogOpts := v1.PodLogOptions{
		Follow: true,
	}

	req := clientset.CoreV1().Pods(service.Namespace).GetLogs(service.PodName, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		stream.errCh <- err
	}
	defer podLogs.Close()

	reader := bufio.NewScanner(podLogs)

	for reader.Scan() {
		line := reader.Text()
		stream.logCh <- line + "\n"
	}
	stream.doneCh <- true
}
