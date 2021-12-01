package utils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

var (
	K8sClient *kubernetes.Clientset = nil
)

func getK8sClient() {
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", "config/k8s-config-home")
	if err != nil {
		log.Printf("K8sUtil.BuildConfigFromFlags: %s\n", err)
		return
	}

	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Printf("K8sUtil.NewForConfig: %s\n", err)
		return
	}

	K8sClient = k8sClient
}

func init() {
	getK8sClient()
}
