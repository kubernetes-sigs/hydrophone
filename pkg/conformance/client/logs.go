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

package client

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"sigs.k8s.io/hydrophone/pkg/conformance"
	"sigs.k8s.io/hydrophone/pkg/log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// Contains all the necessary channels to transfer data
type streamLogs struct {
	logCh  chan string
	errCh  chan error
	doneCh chan bool
}

// PrintE2ELogs checks for Pod and starts a goroutine to print the logs in real-time to stdout.
func (c *Client) PrintE2ELogs(ctx context.Context) error {
	informerFactory := informers.NewSharedInformerFactory(c.clientset, 10*time.Second)

	podInformer := informerFactory.Core().V1().Pods()

	if _, err := podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{}); err != nil {
		return fmt.Errorf("failed to add event handler: %w", err)
	}

	informerFactory.Start(ctx.Done())
	informerFactory.WaitForCacheSync(ctx.Done())

	for {
		pod, _ := podInformer.Lister().Pods(c.namespace).Get(conformance.PodName)
		if pod.Status.Phase == corev1.PodRunning {
			var err error
			stream := streamLogs{
				logCh:  make(chan string),
				errCh:  make(chan error),
				doneCh: make(chan bool),
			}

			go c.streamPodLogs(ctx, stream)

		loop:
			for {
				select {
				case err = <-stream.errCh:
					log.Fatal(err)
				case logStream := <-stream.logCh:
					_, err = fmt.Print(logStream)
					if err != nil {
						log.Fatal(err)
					}
				case <-stream.doneCh:
					break loop
				}
			}
			break
		}
	}

	return nil
}

// List pod resource with the given namespace
func (c *Client) streamPodLogs(ctx context.Context, stream streamLogs) {
	podLogOpts := corev1.PodLogOptions{
		Container: conformance.ConformanceContainer,
		Follow:    true,
	}

	req := c.clientset.CoreV1().Pods(c.namespace).GetLogs(conformance.PodName, &podLogOpts)
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
