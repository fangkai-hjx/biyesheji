package main

import (
	"context"
	"fmt"
	"github.com/TimeBye/go-harbor/pkg/model"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"strconv"
	"strings"
	"t/back/entity"
	"t/back/middleware"
	"t/back/utils"
	"time"
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
		result := make([]entity.Namespace, 0)
		for _, v := range namespaceList.Items {
			var r = entity.Namespace{
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
	// 创建服务
	r.PUT("/pub_service/workspace/:namespace", func(c *gin.Context) {
		namespace := c.Param("namespace")
		service_name := c.PostForm("service_name")
		image := c.PostForm("image")
		//tag := "app=nginx-demo"
		port := strToInt32(c.PostForm("port"))
		//count := c.Param("count")
		count := strToInt32(c.PostForm("count"))
		fmt.Println(namespace, image, service_name, image, port)
		k8sClient := utils.GetK8sClient()
		if k8sClient == nil {
			fmt.Println("k8sClient is nil")
			return
		}
		// check service from db
		dbClient := utils.GetDBClient()
		if dbClient == nil {
			fmt.Println("dbClient is nil")
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    nil,
				"message": "获取数据库连接错误",
			})
			return
		}
		row := dbClient.QueryRow(`SELECT COUNT(*) AS total FROM tb_service WHERE delete_flag!=0 and service_name = ? and namespace= ?`, service_name, namespace)

		var r int64
		if err := row.Scan(&r); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    nil,
				"message": "数据库查询失败",
			})
			return
		}
		if r != 0 {
			c.JSON(http.StatusOK, gin.H{
				"data":    nil,
				"message": "服务已被创建，请勿重复创建",
			})
			return
		}
		// insert or update
		// if
		now := time.Now().Unix()
		row = dbClient.QueryRow(`SELECT COUNT(*) AS total FROM tb_service WHERE service_name = ? and namespace= ?`, service_name, namespace)
		if err := row.Scan(&r); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    nil,
				"message": "数据库查询失败",
			})
			return
		}
		if r == 0 {
			_, err := dbClient.Exec(`INSERT INTO tb_service (service_name,namespace,create_time,replicas,image_name,env_vars,ports,description,creator,change_time) VALUES 
					(?,?,?,?,?,?,?,?,?,?)`,
				service_name, namespace, now, count, image, "ttt", port, "ttttt", "admin", now)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"data":    nil,
					"message": "数据库插入失败" + err.Error(),
				})
			}
			return
		} else {
			_, err := dbClient.Exec(`UPDATE tb_service SET delete_flag=? WHERE service_name = ? and namespace=?`, "1", service_name, namespace)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"data":    nil,
					"message": "数据库更新失败" + err.Error(),
				})
			}
		}
		deploymentsClient := k8sClient.AppsV1().Deployments(namespace)
		servicesClient := k8sClient.CoreV1().Services(namespace)
		service := &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: service_name,
			},
			Spec: apiv1.ServiceSpec{
				Ports: []apiv1.ServicePort{
					apiv1.ServicePort{
						Port: 80,
					},
				},
				Selector: map[string]string{"app": "nginx-demo"},
			},
		}
		deployment := &appsv1.Deployment{ //
			ObjectMeta: metav1.ObjectMeta{
				Name: service_name,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(count),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "nginx-demo",
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx-demo",
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  service_name,
								Image: image,
								Ports: []apiv1.ContainerPort{
									{
										Name:          "http",
										Protocol:      apiv1.ProtocolTCP,
										ContainerPort: port,
									},
								},
							},
						},
					},
				},
			},
		}
		fmt.Println("Creating deployment...")
		dep_result, err1 := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
		svc_result, err2 := servicesClient.Create(context.TODO(), service, metav1.CreateOptions{})
		// 将服务ip 更新到数据库 svc_result.Spec.ClusterIP
		_, err := dbClient.Exec(`UPDATE tb_service SET cluster_ip=? WHERE service_name = ? and namespace= ?`, svc_result.Spec.ClusterIP, service_name,namespace)
		if err != nil {
			fmt.Println("服务IP更新到数据库失败")
			return
		}
		data := ""
		if err1 != nil {
			data += err1.Error()
		}
		if err2 != nil {
			data += err2.Error()
		}
		if data != "" {
			c.JSON(http.StatusOK, gin.H{
				"message": "创建失败",
				"data":    data,
			})
			return
		}
		fmt.Printf("Created deployment %q.\n", dep_result.GetObjectMeta().GetName())
		fmt.Printf("Created service %q.\n", svc_result.GetObjectMeta().GetName())
		c.JSON(http.StatusOK, gin.H{
			"message": "创建成功",
			"data":    nil,
		})
	})
	// 删除服务
	r.DELETE("/pub_service/workspace/:namespace", func(c *gin.Context) {
		namespace := c.Param("namespace")
		service_name := c.PostForm("service_name")
		fmt.Println(namespace, service_name)
		k8sClient := utils.GetK8sClient()
		if k8sClient == nil {
			fmt.Println("k8sClient is nil")
			return
		}
		// check service from db
		dbClient := utils.GetDBClient()
		if dbClient == nil {
			fmt.Println("dbClient is nil")
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    nil,
				"message": "获取数据库连接错误",
			})
			return
		}
		now := time.Now().Unix()
		row := dbClient.QueryRow(`SELECT COUNT(*) AS total FROM tb_service WHERE delete_flag!=0 and service_name = ? and namespace= ?`, service_name,namespace)
		var r int64
		if err := row.Scan(&r); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    nil,
				"message": "数据库查询失败",
			})
			return
		}
		if r == 0 {
			c.JSON(http.StatusOK, gin.H{
				"data":    nil,
				"message": "服务不存在，删除失败",
			})
			return
		}
		// delete
		err := k8sClient.AppsV1().Deployments(namespace).Delete(context.TODO(), service_name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println("delete deployment is error:", err.Error())
			return
		}
		err = k8sClient.CoreV1().Services(namespace).Delete(context.TODO(), service_name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println("delete service is error:", err.Error())
			return
		}
		// 将服务ip 更新到数据库 svc_result.Spec.ClusterIP
		_, err = dbClient.Exec(`UPDATE tb_service SET delete_flag=?,change_time =? WHERE service_name = ? and namespace= ?`, "0", now, service_name,namespace)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"data":    nil,
				"message": "服务删除失败" + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    nil,
			"message": "删除服务" + service_name + "成功",
		})
		return
	})
	// 修改服务参数
		// 1 副本数
		// 2 修改镜像
		// 3 修改端口
	r.POST("/pub_service/workspace/:namespace", func(c *gin.Context) {
		namespace := c.Param("namespace")
		service_name := c.PostForm("service_name")
		replica := c.PostForm("replica")
		//image := c.PostForm("image")
		//port := c.PostForm("port")
		deploymentsClient := utils.GetK8sClient().AppsV1().Deployments(namespace)
		deployment, err := deploymentsClient.Get(context.TODO(), service_name, metav1.GetOptions{})
		if err != nil{
			c.JSON(http.StatusOK,gin.H{
				"data":nil,
				"message":"服务不存在"+err.Error(),
			})
			return
		}
		if deployment != nil{
			deployment.Spec.Replicas = int32Ptr(strToInt32(replica))
		}
		_, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if err != nil{
			c.JSON(http.StatusOK,gin.H{
				"data":nil,
				"message":"服务更新错误",
			})
			return
		}
		dbClient := utils.GetDBClient()
		if dbClient == nil {
			fmt.Println("dbClient is nil")
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    nil,
				"message": "获取数据库连接错误",
			})
			return
		}
		now := time.Now().Unix()
		_, err = dbClient.Exec(`UPDATE tb_service SET replicas=?,change_time =? WHERE service_name = ? and namespace= ?`, replica, now, service_name,namespace)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"data":    nil,
				"message": "服务删除失败" + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"data":nil,
			"message":"服务更新成功",
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

		result := make([]entity.Service, 0)
		for _, v := range serviceList.Items {
			s := entity.Service{
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
func int32Ptr(i int32) *int32 { return &i }
func strToInt32(str string) int32 {
	j, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return -1
	}
	return int32(j)
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
