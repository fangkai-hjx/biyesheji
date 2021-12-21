package image

import (
	"github.com/TimeBye/go-harbor/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"t/back/entity"
	"t/back/utils"
)

// Router 服务发布 镜像管理: 查询所有的项目
func Router(router *gin.RouterGroup) {
	app := router.Group("/img/workspace")
	{
		app.GET("/all", getAllImage)
		app.GET("/project", getImageProject)
		app.GET("/:project_name/img", getImageForProject)
		//app.PUT("/:namespace", addNamespace)
		//app.DELETE("/:namespace", delNamespace)
		//cluster.PUT("/cluster",updateCluster )
		//cluster.PUT("/cluster",updateCluster )
	}
}

// 上传镜像到服务器
func uploadImage(c *gin.Context) {

}
func getImageForProject(c *gin.Context) {
	project_name := c.Param("project_name")
	query := model.Query{}
	HarborClient := utils.GetHarborClient()
	project, err := HarborClient.V2.Get(project_name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "查询服务shibai",
			"data":    nil,
		})
		return
	}
	images, err := HarborClient.V2.Repositories(project.Name).List(&query)
	if images == nil || len(*images) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "查询服务shibai",
			"data":    nil,
		})
		return
	}
	res := make([]entity.Image, 0)
	for _, image := range *images {
		tags, _ := GetSomeImage(strings.Split(image.Name, "/")[0], strings.Split(image.Name, "/")[1])
		for _, v := range tags {
			var img = entity.Image{
				PullCount:    image.PullCount,
				ImageName:    image.Name,
				Tag:          v,
				UpdateTime:   image.UpdateTime.Format("2006-01-02 15:04:05"),
				CreateTime:   image.CreationTime.Format("2006-01-02 15:04:05"),
				ProjectId:    image.ProjectID,
				Description:  image.Description,
				RepositoryId: image.RepositoryID,
			}
			res = append(res, img)
		}

	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询服务shibai",
		"data":    res,
	})
}
func getImageProject(c *gin.Context) {
	query := model.Query{}
	HarborClient := utils.GetHarborClient()
	projects, err := HarborClient.V2.List(&query)
	if err != nil {
		return
	}
	result := make([]string, 0)
	for _, project := range *projects {
		result = append(result, project.Name)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询服务成功",
		"data":    result,
	})
}
func getAllImage(c *gin.Context) {
	query := model.Query{}
	HarborClient := utils.GetHarborClient()
	projects, err := HarborClient.V2.List(&query)
	if err != nil {
		return
	}
	res := make([]entity.HarborRepository, 0)
	for _, project := range *projects {

		images, err := HarborClient.V2.Repositories(project.Name).List(&query)
		if err != nil {

		}
		if images == nil || len(*images) == 0 {
			continue
		}
		var harborRepository = entity.HarborRepository{
			ProjectId:   project.ProjectID,
			ProjectName: project.Name,
			CreateTime:  project.CreationTime.Format("2006-01-02 15:04:05"),
			UpdateTime:  project.UpdateTime.Format("2006-01-02 15:04:05"),
			OwnerName:   project.OwnerName,
			RepoCount:   project.RepoCount,
			Images:      nil,
		}
		//加上镜像
		imageList := make([]entity.Image, 0)
		for _, image := range *images {
			tags, _ := GetSomeImage(strings.Split(image.Name, "/")[0], strings.Split(image.Name, "/")[1])
			for _, v := range tags {
				var img = entity.Image{
					PullCount:    image.PullCount,
					ImageName:    image.Name,
					UpdateTime:   image.UpdateTime.Format("2006-01-02 15:04:05"),
					CreateTime:   image.CreationTime.Format("2006-01-02 15:04:05"),
					ProjectId:    image.ProjectID,
					Description:  image.Description,
					RepositoryId: image.RepositoryID,
					Tag: v,
				}
				imageList = append(imageList, img)
			}
		}
		harborRepository.Images = imageList
		res = append(res, harborRepository)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询服务成功",
		"data":    res,
	})
}
func GetSomeImage(projectName, imageName string) (res []string, err error) {
	var query = model.Query{}
	tags, _ := utils.HarborClient.V2.Repositories(projectName).Artifacts(imageName).List(&query)
	for _, tag := range *tags {
		for _, tag := range tag.Tags {
			res = append(res, tag.Name)
		}
	}
	return res, err
}
