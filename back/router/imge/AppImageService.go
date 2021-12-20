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
		//app.PUT("/:namespace", addNamespace)
		//app.DELETE("/:namespace", delNamespace)
		//cluster.PUT("/cluster",updateCluster )
		//cluster.PUT("/cluster",updateCluster )
	}
}
// 上传镜像到服务器
func uploadImage(c *gin.Context)  {

}
func getAllImage(c *gin.Context)  {
	query := model.Query{}
	HarborClient := utils.GetHarborClient()
	projects, err := HarborClient.V2.List(&query)
	if err != nil {
		return
	}
	res := make([]entity.HarborRepository, 0)
	for _, project := range *projects {
		var harborRepository = entity.HarborRepository{
			ProjectId:   project.ProjectID,
			ProjectName: project.Name,
			CreateTime:  project.CreationTime.String(),
			UpdateTime:  project.UpdateTime.String(),
			OwnerName:   project.OwnerName,
			RepoCount:   project.RepoCount,
			Images:      nil,
		}
		//加上镜像
		images, err := HarborClient.V2.Repositories(project.Name).List(&query)
		if err != nil {

		}
		imageList := make([]entity.Image, 0)
		for _, image := range *images {
			var i = entity.Image{
				PullCount:    image.PullCount,
				ImageName:    image.Name,
				UpdateTime:   image.UpdateTime.String(),
				CreateTime:   image.CreationTime.String(),
				ProjectId:    image.ProjectID,
				Description:  image.Description,
				RepositoryId: image.RepositoryID,
			}
			tags, _ := GetSomeImage(strings.Split(i.ImageName, "/")[0], strings.Split(i.ImageName, "/")[1])
			i.Tag = tags
			imageList = append(imageList, i)
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
