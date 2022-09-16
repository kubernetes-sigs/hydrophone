package client

import (
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

const (
	NAMESPACE = "conformance"
	POD_NAME  = "e2e-conformance-test"
)

// Check for Pod and start a go routine if new deployment added
func (c *Client) CheckForE2ELogs() {
	labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
	})

	informerFactory := informers.NewSharedInformerFactoryWithOptions(c.ClientSet, 2*time.Minute,
		informers.WithNamespace(NAMESPACE), labelOptions)

	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Println("Pod added. Let's start checking!")

			ch := make(chan error)
			done := make(chan bool)

			go c.getLogs(ch, done)

		loop:
			for {
				select {
				case err := <-ch:
					log.Fatalf("error getting logs: %v", err)
				case <-done:
					break loop
				}
			}
		},
	})

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
}

func (c *Client) getLogs(ch chan error, done chan bool) {
	err := GetPodLogs(c.ClientSet)
	if err != nil {
		ch <- fmt.Errorf("get logs: %s", err.Error())
	}

	done <- true
}
