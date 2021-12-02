package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"t/back/constant"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	JwtInstance *JWT
	//JwtClaim    *jwt.StandardClaims
)

func init() {
	JwtInstance = newJWT()
	//JwtClaim = newJwtClaim()
}

func IsTokenStringValid(str *string) bool {
	*str = strings.Replace(*str, " ", "", -1)
	// 去除换行符
	*str = strings.Replace(*str, "\n", "", -1)
	return len(*str) > 0
}

/*JWTAuth 中间件，检查token*/
func JWTUserTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("token")
		if !IsTokenStringValid(&token) {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}

		log.Print("get token: ", token)

		j := newJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseUserToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    "授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)

	}
}

func UserTokenAuthMid() gin.HandlerFunc {
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

		//获得token
		token := c.Request.Header.Get("Authorization")
		log.Printf("the token get from request is :%s", token)

		//判断token是否为空
		if !IsTokenStringValid(&token) {
			log.Println("token not valid")
			c.Writer.Header().Set("tokenCode", fmt.Sprintf("%d", constant.TokenNotValid))
			//c.Header("tokenCode", fmt.Sprintf("%d", constant.TokenValid))
			c.JSON(http.StatusOK, gin.H{
				"returnCode":  400,
				"description": "token不存在",
			})
			c.Abort()
			return
			//fmt.Errorf("The token is not valid :%s ", token)
		}

		//解析token获得信息
		if claims, err := JwtInstance.ParseUserToken(token); err != nil {
			if err == TokenExpired {
				log.Println("UserTokenAuth()-token expire")
				c.Writer.Header().Set("tokenCode", fmt.Sprintf("%d", constant.TokenExpireCode))
				//c.Header("tokenCode", fmt.Sprintf("%d", constant.TokenValid))
				c.JSON(http.StatusOK, gin.H{
					"returnCode":  400,
					"description": "token过期",
				})
				c.Abort()

				//return constant.TokenExpireCode, err
			} else {
				log.Println("UserTokenAuth()-token error for some reasons")
				c.Writer.Header().Set("tokenCode", fmt.Sprintf("%d", constant.TokenNotValid))
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
			c.Writer.Header().Set("tokenCode", fmt.Sprintf("%d", constant.TokenValid))
			//c.Header("tokenCode", fmt.Sprintf("%d", constant.TokenValid))
			c.Next()

		}
		return

	}
}

/*检查token是否有效的函数,返回返回码以及错误*/
func UserTokenAuth(c *gin.Context) (returnCode int, err error) {

	//白名单内的路径不需要验证token
	requestUrl := c.Request.RequestURI
	for _, url := range constant.WhiteUrlList {
		if url == requestUrl {
			log.Println("不验证token url in white list")
			return constant.TokenValid, nil
		}
	}

	//获得token
	token := c.Request.Header.Get("Authorization")
	url := c.Request.RequestURI
	log.Println("request url is:" + url)
	log.Printf("the token get from request is :%s", token)

	//判断token是否为空
	if !IsTokenStringValid(&token) {
		log.Println("token not valid")
		return constant.TokenNotValid, fmt.Errorf("The token is not valid :%s ", token)
	}

	//解析token获得信息
	if claims, err := JwtInstance.ParseUserToken(token); err != nil {
		if err == TokenExpired {
			log.Println("UserTokenAuth()-token expire")
			return constant.TokenExpireCode, err
		} else {
			log.Println("UserTokenAuth()-token error for some reasons")
			return constant.TokenNotValid, err
		}
	} else {
		c.Set("claims", claims)
		log.Printf("UserTokenAuth Pass,the clasims are %v", claims)
		return constant.TokenValid, nil
	}

}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 一些常量
var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
	SignKey          = "ServicePlatForm"
	Issuer           = "admin"
)

//type CustomClaims struct {
//	ID    string `json:"userId"`
//	Name  string `json:"name"`
//	Phone string `json:"phone"`
//	jwt.StandardClaims
//}

// 载荷，可以加一些自己需要的信息
type UserInfoClaims struct {
	UserName string `json:"userName"`
	Power    int    `json:"power"`
	Email    string `json:"phone"`
	jwt.StandardClaims
}

// 获取signKey
func getSignKey() string {
	return SignKey
}

// 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

// CreateToken 生成一个token
func (j *JWT) CreateUserToken(claims UserInfoClaims) (string, error) {
	//claims :=  UserInfoClaims{UserName: userInfo.UserName,Power:userInfo.Power,Email: userInfo.Email}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

/*解析token,返回详情以及错误*/
func (j *JWT) ParseUserToken(tokenString string) (*UserInfoClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserInfoClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*UserInfoClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func (j *JWT) RefreshUserToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserInfoClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*UserInfoClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateUserToken(*claims)
	}
	return "", TokenInvalid
}

// 新建一个jwt实例
func newJWT() *JWT {
	return &JWT{
		[]byte(getSignKey()),
	}
}

//新建一个 StandardClaims实例
func NewJwtClaim() *jwt.StandardClaims {

	return &jwt.StandardClaims{
		NotBefore: time.Now().Unix() - constant.TokenBefore, // 签名生效时间
		ExpiresAt: time.Now().Unix() + constant.TokenExpire, // 过期时间 一小时
		Issuer:    Issuer,                                   //签名的发行者
	}

}

//// 更新token
//func (j *JWT) RefreshToken(tokenString string) (string, error) {
//	jwt.TimeFunc = func() time.Time {
//		return time.Unix(0, 0)
//	}
//	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return j.SigningKey, nil
//	})
//	if err != nil {
//		return "", err
//	}
//	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
//		jwt.TimeFunc = time.Now
//		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
//		return j.CreateToken(*claims)
//	}
//	return "", TokenInvalid
//}
//// CreateToken 生成一个token
//func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	return token.SignedString(j.SigningKey)
//}
//
//// 解析Token
//func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
//	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return j.SigningKey, nil
//	})
//	if err != nil {
//		if ve, ok := err.(*jwt.ValidationError); ok {
//			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
//				return nil, TokenMalformed
//			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
//				// Token is expired
//				return nil, TokenExpired
//			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
//				return nil, TokenNotValidYet
//			} else {
//				return nil, TokenInvalid
//			}
//		}
//	}
//	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
//		return claims, nil
//	}
//	return nil, TokenInvalid
//}

//// JWTAuth 中间件，检查token
//func JWTAuth() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		token := c.Request.Header.Get("token")
//		if token == "" {
//			c.JSON(http.StatusOK, gin.H{
//				"status": -1,
//				"msg":    "请求未携带token，无权限访问",
//			})
//			c.Abort()
//			return
//		}
//
//		log.Print("get token: ", token)
//
//		j := newJWT()
//		// parseToken 解析token包含的信息
//		claims, err := j.ParseToken(token)
//		if err != nil {
//			if err == TokenExpired {
//				c.JSON(http.StatusOK, gin.H{
//					"status": -1,
//					"msg":    "授权已过期",
//				})
//				c.Abort()
//				return
//			}
//			c.JSON(http.StatusOK, gin.H{
//				"status": -1,
//				"msg":    err.Error(),
//			})
//			c.Abort()
//			return
//		}
//		// 继续交由下一个路由处理,并将解析出的信息传递下去
//		c.Set("claims", claims)
//	}
//}
