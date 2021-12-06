package utils

import (
	"github.com/TimeBye/go-harbor"
	client2 "github.com/TimeBye/go-harbor/pkg/client"
	"log"
)

var (
	HarborClient *client2.Clientset = nil
)

func GetHarborClient() *client2.Clientset {

	if HarborClient != nil {
		return HarborClient
	}
	harborClient, err := harbor.NewClientSet("http://116.56.140.62:8089/", "root", "Harbor12345")
	if err != nil {
		log.Println("GetHarborClient error : ", err)
		return nil
	}

	HarborClient = harborClient
	return HarborClient
}
