package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"strings"
	_ "t/back/config"
	"t/back/utils"
)

func main() {
	err := utils.InitClient()
	if err != nil {
		fmt.Println("init redis failed..")
	}

	k8sConfig, err := clientcmd.BuildConfigFromFlags("", "k8s-config")
	if err != nil {
		log.Printf("K8sUtil.BuildConfigFromFlags: %s\n", err)
		return
	}
	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Printf("K8sUtil.NewForConfig: %s\n", err)
		return
	}
	namespace := "hpa"
	w, _ := k8sClient.CoreV1().Events(namespace).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("create watch error, error is %s, program exit!", err.Error())
		panic(err)
	}
	redisClient := utils.Rdb
	var key string = "event:warning"
	for {
		select {
		case <-w.ResultChan():
			www := <-w.ResultChan()
			if www.Object != nil {
				//panic: interface conversion: runtime.Object is *v1.Status, not *v1.Even
				t := www.Object.(*v1.Event)
				//fmt.Println(t.Type, t.Name, t.EventTime, t.Reason)
				//	////Warning nginx-6d65fc45c6-thvsl.16b0a0d65e45a10d 0001-01-01 00:00:00 +0000 UTC Failed
				if t.Type == "Warning" {
					var str = strings.Split(t.Name, ".")[0]
					name := GetHostName(k8sClient, namespace, str)
					fmt.Printf("node[%s] happen warning event[%s]\n",name,t.Reason)
					////发送消息给消息队列
					redisClient.LPush(key, name)
				}
			}

		}
	}
}
func GetHostName(k8sClient *kubernetes.Clientset, namespace, name string) string {
	fmt.Println(name)
	w, _ := k8sClient.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return w.Spec.NodeName
}
