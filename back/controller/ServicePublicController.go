package controller

import (
	"gee/entity"
	"gee/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ServicePublicController struct {
	service       *service.ApplicationService
}

/* 创建 ServiceController 实例 */
func NewServiceController() *ServicePublicController {
	controller := &ServicePublicController{
	}
	return controller
}

func (controller *ServicePublicController) PublishService(context *gin.Context) {
	var request entity.PublishServiceRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		log.Printf("ServiceController.PublishService.ShouldBindJSON: %s\n", err)

		context.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "data bind failed",
			"data":    err,
		})
		return
	}

	if !request.IsValid() {
		context.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "invalid request",
			"data":    "",
		})
		return
	}

	if !controller.service.CheckResourceInventory(&request.Service) {
		context.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "resource is insufficient",
			"data":    "",
		})
		return
	}

	address, err := controller.service.GetImageRepository(request.Image)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "publish service failed",
			"data":    err,
		})
		return
	}
	request.Image = address + "/" + request.Image
	err = controller.service.PublishService(&request.Service)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "publish service failed",
			"data":    err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "publish service successfully",
		"data":    request.Service,
	})
}