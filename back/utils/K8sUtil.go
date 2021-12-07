package utils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

var (
	K8sClient *kubernetes.Clientset = nil
)

func GetK8sClient() *kubernetes.Clientset{
	if K8sClient != nil{
		return K8sClient
	}
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", "../back/config/k8s-config")
	if err != nil {
		log.Printf("K8sUtil.BuildConfigFromFlags: %s\n", err)
		return nil
	}

	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Printf("K8sUtil.NewForConfig: %s\n", err)
		return nil
	}

	K8sClient = k8sClient
	return K8sClient
}