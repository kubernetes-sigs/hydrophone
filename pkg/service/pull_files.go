package service

import (
	"context"

	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/hydrophone/pkg/common"
	"sigs.k8s.io/hydrophone/pkg/log"
)

func PullFiles(clientset *kubernetes.Clientset) {

	namespace := viper.GetString("namespace")

	retrievalPod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      common.RetrievalPodName,
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    common.OutputContainer,
					Image:   viper.GetString("busybox-image"),
					Command: []string{"/bin/sh", "-c", "sleep infinity"},
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
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: common.PVCName,
						},
					},
				},
			},
			RestartPolicy:      v1.RestartPolicyNever,
			ServiceAccountName: "conformance-serviceaccount",
		},
	}

	pod, err := clientset.CoreV1().Pods(namespace).Create(ctx, &retrievalPod, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Fatalf("pod already exist %s. Please run cleanup first", retrievalPod.ObjectMeta.Name)
		} else {
			log.Fatal(err)
		}
	}
	log.Printf("pod created %s\n", pod.Name)

	// Wait for the pod to be running
	watchOptions := metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("metadata.name", common.RetrievalPodName).String(),
	}

	watchInterface, err := clientset.CoreV1().Pods(namespace).Watch(context.Background(), watchOptions)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Waiting for pod to be running and ready...")
	for event := range watchInterface.ResultChan() {
		pod, ok := event.Object.(*v1.Pod)
		if !ok {
			continue
		}

		isRunningAndReady := pod.Status.Phase == v1.PodRunning
		for _, cond := range pod.Status.Conditions {
			if cond.Type == v1.PodReady && cond.Status == v1.ConditionTrue {
				isRunningAndReady = isRunningAndReady && true
				break
			}
		}

		if isRunningAndReady {
			break
		}
	}
}
