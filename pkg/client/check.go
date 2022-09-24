package client

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dims/k8s-run-e2e/pkg/service"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// Check for Pod and start a go routine if new deployment added
func (c *Client) CheckForE2ELogs() {
	informerFactory := informers.NewSharedInformerFactory(c.ClientSet, 10*time.Second)

	podInformer := informerFactory.Core().V1().Pods()

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)

	for {
		pod, _ := podInformer.Lister().Pods(service.Namespace).Get(service.PodName)
		if pod.Status.Phase == v1.PodRunning {
			c.getLogs()
			break
		}
	}
}

func (c *Client) getLogs() {
	cancelCtx := context.Background()
	cancelCtx, cancelFunc := context.WithCancel(cancelCtx)
	defer cancelFunc()

	err := getPodLogs(cancelCtx, c.ClientSet)
	if err != nil {
		log.Fatal(err)
	}
}
