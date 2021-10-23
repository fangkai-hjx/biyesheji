package main

import (
	"github.com/gok8s/kube-eventalert/pkg/config"
	"github.com/gok8s/kube-eventalert/pkg/controller"
	"github.com/gok8s/kube-eventalert/pkg/store"
	api_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", "/root/my-sync-job/config/k8s-config")
	if err != nil {
		log.Printf("K8sUtil.BuildConfigFromFlags: %s\n", err)
		return
	}
	var config config.Config
	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Printf("K8sUtil.NewForConfig: %s\n", err)
		return
	}
	lw := cache.NewListWatchFromClient(
		k8sClient.CoreV1().RESTClient(), // 客户端
		"events",                         // 被监控资源类型
		"",                               // 被监控命名空间
		fields.Everything())
	informer := cache.NewSharedIndexInformer(lw, &api_v1.Event{}, 0, cache.Indexers{})
	c := controller.NewResourceController(k8sClient, informer, config)
	c.MQClient, err = store.NewRabbitMQClient(config)

}
