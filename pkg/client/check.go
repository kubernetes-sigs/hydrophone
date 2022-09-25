package client

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dims/k8s-run-e2e/pkg/service"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// Contains all the necessary channels to transfer data
type streamLogs struct {
	logCh  chan string
	errCh  chan error
	doneCh chan bool
}

// Check for Pod and start a go routine if new deployment added
func (c *Client) CheckForE2ELogs(output string) {
	informerFactory := informers.NewSharedInformerFactory(c.ClientSet, 10*time.Second)

	podInformer := informerFactory.Core().V1().Pods()

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)

	for {
		pod, _ := podInformer.Lister().Pods(service.Namespace).Get(service.PodName)
		if pod.Status.Phase == v1.PodRunning {
			file, err := createLogDirAndFile(output)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			stream := streamLogs{
				logCh:  make(chan string),
				errCh:  make(chan error),
				doneCh: make(chan bool),
			}

			go getPodLogs(c.ClientSet, stream)
			if err != nil {
				log.Fatal(err)
			}

		loop:
			for {
				select {
				case err = <-stream.errCh:
					log.Fatal(err)
				case logStream := <-stream.logCh:
					_, err = file.WriteString(logStream)
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
}

// createLogDirAndFile create a directory and create a file of name same as pod
func createLogDirAndFile(output string) (*os.File, error) {
	if err := os.Mkdir(output, os.ModePerm); err != nil {
		return nil, err
	}

	path := filepath.Join(output, service.PodName)

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
