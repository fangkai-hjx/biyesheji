package main

import (
	"context"
	"fmt"
	"github.com/TimeBye/go-harbor/pkg/model"
	"github.com/gin-gonic/gin"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"strings"
	"t/back/middleware"
	"t/back/utils"
)

func main() {
	r := gin.Default()
	r.POST("/login", middleware.LoginCheck(), Login)
	r.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"dada":  400,
			"dadsa": "dadad",
		})
	})
	// 服务发布
	// 镜像管理
	// 查询所有的项目
	r.GET("/image/project/all", func(c *gin.Context) {
		query := model.Query{}
		HarborClient := utils.GetHarborClient()
		projects, err := HarborClient.V2.List(&query)
		if err != nil {
			return
		}
		res := make([]HarborRepository, 0)
		for _, project := range *projects {
			var harborRepository = HarborRepository{
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
			imageList := make([]Image, 0)
			for _, image := range *images {
				var i = Image{
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
	})
	// 查询项目下所有的镜像
	//r.GET("/image/project/images", func(c *gin.Context) {
	//	harborClient := utils.GetHarborClient()
	//	if harborClient == nil {
	//
	//		fmt.Println("harborClient is nil")
	//		return
	//	}
	//	query := model.Query{}
	//	repositories, _ := harborClient.V2.Repositories("public").List(&query)
	//	for _, p := range *repositories {
	//		p.
	//	}
	//})

	// 服务发布
	// 工作空间 管理
	// 查询所有的工作空间 http://localhost:8080/pub_service/workspace/all
	r.GET("/pub_service/workspace/all", func(c *gin.Context) {
		client := utils.GetK8sClient()
		if client == nil {
			fmt.Println("client is nill")
			return
		}
		namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
		result := make([]Namespace, 0)
		for _, v := range namespaceList.Items {
			var r = Namespace{
				Name:   v.Name,
				Status: string(v.Status.Phase),
			}
			result = append(result, r)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "查询成功",
			"data":    result,
		})
	})
	// 创建指定的工作空间
	r.GET("/pub_service/workspace/:namespace", func(c *gin.Context) {
		//namespace := c.PostForm("namespace")
		namespace := c.Param("namespace")
		client := utils.GetK8sClient()
		if client == nil {
			fmt.Println("client is nill")
			return
		}
		ns := &apiv1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:   namespace,
				Labels: map[string]string{"name": namespace},
			},
		}
		_, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "增加失败",
				"data":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "增加" + namespace + "成功",
			"data":    nil,
		})
		return
	})
	// 删除指定的工作空间
	r.GET("/pub_service/workspace/del/:namespace", func(c *gin.Context) {
		//namespace := c.PostForm("namespace")
		namespace := c.Param("namespace")
		client := utils.GetK8sClient()
		if client == nil {
			fmt.Println("client is nill")
			return
		}
		err := client.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "删除失败",
				"data":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "删除" + namespace + "成功",
			"data":    nil,
		})
	})

	// 服务相关
	r.GET("/pub_service/:namespace/all", func(c *gin.Context) {
		namespace := c.Param("namespace")
		client := utils.GetK8sClient()
		if client == nil {
			fmt.Println("client is nill")
			return
		}
		serviceList, err := client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "查询服务失败",
				"data":    err.Error(),
			})
			return
		}
		type Service struct {
			Name            string `json:"name"`
			ClusterIP       string `json:"cluster_ip"`
			SessionAffinity string `json:"session_affinity"`
			Status          string `json:"status"`
		}
		result := make([]Service, 0)
		for _, v := range serviceList.Items {
			s := Service{
				Name:            v.Name,
				ClusterIP:       v.Spec.ClusterIP,
				SessionAffinity: string(v.Spec.SessionAffinity),
				Status:          v.Status.String(),
			}
			result = append(result, s)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "查询服务成功",
		})
	})
	r.GET("/pub_service/:namespace/:svc_name")
	// 服务治理
	// 自动运维
	r.GET("/service_log/:namespace/:svc_name")
	r.GET("/pod/:namespace/:pod_name")
	r.GET("/pod_log/:namespace/:pod_name")
	r.Run() // listen and serve on 0.0.0.0:8080
}

// 登录成功-->发送token-->跳转到主页
// 登录失败-->调整到登录页
//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6ImFkbWluIiwicG93ZXIiOjAsInBob25lIjoiIiwiZXhwIjoxNjM4MzMwNDUwLCJpc3MiOiJhZG1pbiIsIm5iZiI6MTYzODMyNzY1MH0.JF9RYXvyZy5Cdzpp89oF6IMLq95gLx7023DIf7wXtmo func
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "admin" && password == "111" {
		//生成jwt的token
		token, _ := generateToken(username, password)
		c.Request.Header.Set("Authorization", token)
		c.Redirect(http.StatusMovedPermanently, "/admin")
	}
}

// 生成令牌
func generateToken(username, password string) (string, error) {
	j := utils.JwtInstance
	claims := utils.UserInfoClaims{
		UserName:       username,
		StandardClaims: *utils.NewJwtClaim(),
	}
	token, err := j.CreateUserToken(claims)
	if err != nil {
		return "", err
	}
	log.Printf("generateToken() get the token,it is : %s", token)
	return token, nil
}

type Namespace struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
type RepositoryAddr struct {
	ImageName   string `json:"imageName"`
	ProjectName string `json:"projectName"`
}
type HarborRepository struct {
	ProjectId   int64  `json:"project_id"`
	ProjectName string `json:"project_name"`
	CreateTime  string `json:"creation_time"`
	UpdateTime  string `json:"create_time"`
	OwnerName   string `json:"owner_name"`
	RepoCount   int64  `json:repo_count`
	Images      []Image
}

//查询某个镜像的Tag
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

type Image struct {
	ImageName    string   `json:"image_name"`
	UpdateTime   string   `json:"update_time"`
	CreateTime   string   `json:"create_time"`
	ProjectId    int64    `json:"project_id"`
	Description  string   `json:"description"`
	RepositoryId int64    `json:"repository_id"`
	PullCount    int64    `json:"pull_count"`
	Tag          []string `json:"tag"`
}
