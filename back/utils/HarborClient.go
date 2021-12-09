package utils

import (
	"github.com/TimeBye/go-harbor"
	client2 "github.com/TimeBye/go-harbor/pkg/client"
	"log"
	"t/back/config"
)

var (
	HarborClient *client2.Clientset = nil
)

func GetHarborClient() *client2.Clientset {

	if HarborClient != nil {
		return HarborClient
	}
	harborClient, err := harbor.NewClientSet(config.ProjectConfig.HarborUrl, config.ProjectConfig.HarborUsername, config.ProjectConfig.HarborPassword)
	if err != nil {
		log.Println("GetHarborClient error : ", err)
		return nil
	}

	HarborClient = harborClient
	return HarborClient
}
