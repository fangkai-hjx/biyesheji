package middleware

import (
	"errors"
	"fmt"
	"t/back/constant"
	"t/back/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)
var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
	SignKey          = "ServicePlatForm"
	Issuer           = "admin"
)

func LoginCheck() gin.HandlerFunc  {
	return func(c *gin.Context) {
		//白名单内的路径不需要验证token
		requestUrl := c.Request.RequestURI
		log.Println("request url is:" + requestUrl)
		for _, url := range constant.WhiteUrlList {
			if url == requestUrl {
				log.Println("不验证token url in white list")
				c.Writer.Header().Set("tokenCode", fmt.Sprintf("%d", constant.TokenValid))
				//c.Header("tokenCode", fmt.Sprintf("%d", constant.TokenValid))
				c.Next()
				return
			}
		}
		token := c.Request.Header.Get("Authorization")
		log.Printf("the token get from request is :%s", token)
		//判断token是否为空
		if !IsTokenStringValid(&token) {
			log.Println("token not valid")
			c.Writer.Header().Set("tokenCode", fmt.Sprintf("%d", TokenNotValid))
			//c.Header("tokenCode", fmt.Sprintf("%d", constant.TokenValid))
			c.JSON(http.StatusOK, gin.H{
				"returnCode":  400,
				"description": "token不存在",
			})
			c.Abort()
			return
		}
		//解析token获得信息
		if claims, err := utils.JwtInstance.ParseUserToken(token); err != nil {
			if err == TokenExpired {
				log.Println("UserTokenAuth()-token expire")
				c.Writer.Header().Set("tokenCode", fmt.Sprintf("%d", TokenExpireCode))
				//c.Header("tokenCode", fmt.Sprintf("%d", constant.TokenValid))
				c.JSON(http.StatusOK, gin.H{
					"returnCode":  400,
					"description": "token过期",
				})
				c.Abort()

				//return constant.TokenExpireCode, err
			} else {
				log.Println("UserTokenAuth()-token error for some reasons")
				c.Writer.Header().Set("tokenCode", fmt.Sprintf("%d", TokenNotValid))
				//c.Header("tokenCode", fmt.Sprintf("%d", constant.TokenValid))
				c.JSON(http.StatusOK, gin.H{
					"returnCode":  400,
					"description": "token验证不通过",
				})
				c.Abort()
			}
		} else {
			c.Set("claims", claims)
			log.Printf("UserTokenAuth Pass,the clasims are %v", claims)
			c.Writer.Header().Set("tokenCode", fmt.Sprintf("%d", TokenValid))
			c.Next()
		}
		return
	}
}
func IsTokenStringValid(str *string) bool {
	*str = strings.Replace(*str, " ", "", -1)
	// 去除换行符
	*str = strings.Replace(*str, "\n", "", -1)
	return len(*str) > 0
}
const (
	TokenExpireCode = 441
	TokenValid      = 241
	TokenNotValid   = 451
	TokenBefore = 1000
	TokenExpire       = 1800
	EmptyString       = ""
	DefaultHarborAddr = "reg.image.com"
)