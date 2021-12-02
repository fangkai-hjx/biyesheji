package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
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
	r.GET("/service/all")
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
		fmt.Println("================")
		result := make([]string, 0)
		for _, v := range namespaceList.Items {
			result = append(result, v.Name)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "查询成功",
			"data":    result,
		})
	})
	// 创建指定的工作空间
	r.POST("/pub_service/workspace", func(c *gin.Context) {
		namespace := c.PostForm("namespace")
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
			"message": "删除成功",
			"data":    nil,
		})
	})
	// 删除指定的工作空间
	r.DELETE("/pub_service/workspace", func(c *gin.Context) {
		namespace := c.PostForm("namespace")
		// 创建命名空间
		fmt.Println("删除工作空间", namespace)
	})
	r.GET("/pub_service/:namespace/all")
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
