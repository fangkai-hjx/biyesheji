package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Router 服务发布
func Router(router *gin.RouterGroup) {
	app := router.Group("/system/config/")
	{
		app.GET("/menu", getSystemConfig)
	}
}

type SysMenu struct {
	Id          int       `json:"id"`
	Url         string    `json:"url"` // 做权限处理
	Path        string    `json:"path"`
	Component   string    `json:"component"`
	Name        string    `json:"name"`
	IconCls     string    `json:"icon_cls"`
	KeepAlive   string    `json:"keep_alive"`
	RequireAuth string    `json:"require_auth"`
	ParentId    int       `json:"parent_id"`
	Enabled     string    `json:"enabled"`
	Children    []SysMenu `json:"children"`
}

func getSystemConfig(c *gin.Context) {
	sm := SysMenu{
		Id:          2,
		Url:         "/",
		Path:        "/home",
		Component:   "Home",
		Name:        "服务发布",
		IconCls:     "fa fa-user-circle-o",
		KeepAlive:   "true",
		RequireAuth: "",
		ParentId:    1,
		Enabled:     "true",
		Children:    []SysMenu{
			{
				Id:          1,
				Url:         "/pub_service/image",
				Path:        "/img",
				Component:   "ImageManager",
				Name:        "镜像管理",
				IconCls:     "",
				KeepAlive:   "true",
				RequireAuth: "",
				ParentId:    2,
				Enabled:     "true",
				Children:    nil,
			},
			{
				Id:          2,
				Url:         "",
				Path:        "/svc",
				Component:   "ServicePublish",
				Name:        "服务发布",
				IconCls:     "",
				KeepAlive:   "true",
				RequireAuth: "",
				ParentId:    2,
				Enabled:     "",
				Children:    nil,
			},
		},
	}
	c.JSON(http.StatusOK,sm)
}
