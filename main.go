package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	utils2 "t/utils"
)

func main() {
	err := utils2.InitClient()
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
	redisClient := utils2.Rdb
	var message string
	var key string = "event:warning"
	for {
		select {
		case <-w.ResultChan():
			www := <-w.ResultChan()
			if www.Object != nil {
				t := www.Object.(*v1.Event)
				fmt.Println(t.Type, t.Name, t.EventTime, t.Reason)
				if t.Type == "Warning" {
					fmt.Println("======出现了错误事件======")
					//发送消息给消息队列
					message = fmt.Sprintf("%s:%s:%s", t.Type, t.Name, t.Reason)
					//Warning nginx-6d65fc45c6-thvsl.16b0a0d65e45a10d 0001-01-01 00:00:00 +0000 UTC Failed
					redisClient.LPush(key, message)
				}
			}

		}
	}
}
