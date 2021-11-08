package main

//import (
//	"fmt"
//	docker "github.com/fsouza/go-dockerclient"
//)
//
//func main() {
//	//client, err := docker.NewClientFromEnv()
//	client, err := docker.NewClient("tcp://222.201.187.166:2375")
//	if err != nil {
//		panic(err)
//	}
//	client.CreateContainer(docker.CreateContainerOptions{
//		Name:             "",
//		Config:           nil,
//		HostConfig:       nil,
//		NetworkingConfig: nil,
//		Context:          nil,
//	})
//	// 前端传送
//	// 训练环境  比如tensorflow
//	// 训练版本  比如 2.0
//	//tensorflow/tensorflow   latest-py3-jupyter
//	//imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
//	//if err != nil {
//	//	panic(err)
//	//}
//	//for _, img := range imgs {
//	//	fmt.Println("ID: ", img.ID)
//	//	fmt.Println("RepoTags: ", img.RepoTags)
//	//	fmt.Println("Created: ", img.Created)
//	//	fmt.Println("Size: ", img.Size)
//	//	fmt.Println("VirtualSize: ", img.VirtualSize)
//	//	fmt.Println("ParentId: ", img.ParentID)
//	//}
//	containers, err := client.ListContainers(docker.ListContainersOptions{})
//	if err != nil{
//		panic(err)
//	}
//	for _, container := range containers {
//		fmt.Println(container.Image)
//		fmt.Println(container.Status)
//	}
//	//使用docker 来隔离运行时环境
//	container, err := client.CreateContainer(docker.CreateContainerOptions{
//		Name:             "tf:v2",
//		Config:           &docker.Config{
//			Image:  "tensorflow/tensorflow:latest-py3-jupyter",
//		},
//		HostConfig:       nil,
//		NetworkingConfig: nil,
//		Context:          nil,
//	})
//	fmt.Println(container.Image)
//}