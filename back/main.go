package main

import (
	"github.com/TimeBye/go-harbor/pkg/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"t/back/config"
	image "t/back/router/imge"
	namespace "t/back/router/namespace"
	service "t/back/router/service"
	"t/back/router/system"
	"t/back/utils"
)

var Router *gin.Engine

func main() {
	Router = gin.Default()
	setRouter()
	Router.Run(":" + config.ProjectConfig.ProjectPort)
}

func setRouter() {
	apiv2 := Router.Group("/api/v1")
	rest := apiv2.Group("/rest")
	namespace.Router(rest)
	service.Router(rest)
	image.Router(rest)
	system.Router(rest)
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
