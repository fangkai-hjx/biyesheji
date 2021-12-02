package main

import (
	"github.com/gin-gonic/gin"
	"github.com/afocus/captcha"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"t/back/middleware"
	"t/utils"
)

func main() {
	r := gin.Default()
	cap := captcha.New()
	// 设置字体
	cap.SetFont("./back/comic.ttf")
	cap.SetSize(64, 32)
	cap.SetDisturbance(captcha.MEDIUM)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	// 获取验证码
	r.GET("/captcha", func(c *gin.Context) {

		// 创建验证码 4个字符 captcha.NUM 字符模式数字类型
		// 返回验证码图像对象以及验证码字符串 后期可以对字符串进行对比 判断验证
		img, _ := cap.Create(4, captcha.NUM)
		_ = png.Encode(c.Writer, img)
	})
	r.POST("/login",middleware.LoginCheck(), Login)
	r.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"dada":  400,
			"dadsa": "dadad",
		})
	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
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