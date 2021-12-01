package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"t/back/middleware"
	"t/back/utils"
)

func main() {
	r := gin.Default()
	r.POST("/login",middleware.LoginCheck(), Login)
	r.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"dada":  400,
			"dadsa": "dadad",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
// 登录成功-->发送token-->跳转到主页
// 登录失败-->调整到登录页
//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6ImFkbWluIiwicG93ZXIiOjAsInBob25lIjoiIiwiZXhwIjoxNjM4MzMwNDUwLCJpc3MiOiJhZG1pbiIsIm5iZiI6MTYzODMyNzY1MH0.JF9RYXvyZy5Cdzpp89oF6IMLq95gLx7023DIf7wXtmo
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "admin" && password == "111"{
		//生成jwt的token
		token, _ := generateToken(username, password)
		c.Request.Header.Set("Authorization",token)
		c.Redirect(http.StatusMovedPermanently,"/admin")
	}
}
// 生成令牌
func generateToken(username,password string) (string, error) {
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