package router

import (
	"github.com/gin-gonic/gin"
)

var (
	Router = gin.Default()

	//// 路由分组
	ServiceGroup    *gin.RouterGroup
	//FlavorGroup     *gin.RouterGroup
	//SourceCodeGroup *gin.RouterGroup
	//ProjectGroup    *gin.RouterGroup
	//UserGroup       *gin.RouterGroup
	//ImageGroup      *gin.RouterGroup
	//
	//// 控制器
	//ServiceController    = controller.NewServiceController()
	//FlavorController     = controller.NewFlavorController()
	//UserController       = controller.NewUserController()
	//ImageController      = controller.NewImageController()
	//SourceCodeController = controller.NewSourceCodeController()
	//ProjectController    = controller.NewProjectController()
)

func init() {

	//Router.Use(middleware.Cors())
	////Router.Use(middleware.UserTokenAuthMid())
	ServiceGroup = Router.Group("/service")
	//FlavorGroup = Router.Group("/flavor")
	//SourceCodeGroup = Router.Group("/source-code")
	//ProjectGroup = Router.Group("/project")
	//UserGroup = Router.Group("/user")
	//ImageGroup = Router.Group("/image")
	//// 服务路由
	ServiceGroup.POST("/publish", ServiceController.PublishService)
	//ServiceGroup.POST("/list", ServiceController.ListServices)
	//ServiceGroup.POST("/show", ServiceController.GetServiceDetail)
	//ServiceGroup.POST("/update", ServiceController.UpdateService)
	//ServiceGroup.POST("/delete", ServiceController.DeleteServices)
	//
	//FlavorGroup.POST("/create", FlavorController.CreateFlavor)
	//FlavorGroup.POST("/list", FlavorController.ListFlavors)
	//FlavorGroup.POST("/show", FlavorController.GetFlavorDetail)
	//FlavorGroup.POST("/delete", FlavorController.DeleteFlavor)
	//
	////获取项目列表信息
	//SourceCodeGroup.GET("/project-list", SourceCodeController.GetProjects)
	////获取 某个项目 的 全部分支
	//SourceCodeGroup.GET("/source-code-branch", SourceCodeController.GetBranchByPojId)
	////创建构建任务
	//SourceCodeGroup.GET("/get-code-packaged", SourceCodeController.GetBuildById)
	//SourceCodeGroup.POST("/create-buildjob", SourceCodeController.CreateBuild)
	//ProjectGroup.POST("/list", ProjectController.GetProjectList)
	//ProjectGroup.POST("/create-project", ProjectController.CreateNewProject)
	//ProjectGroup.POST("/edit-project", ProjectController.EditProject)
	//ProjectGroup.POST("/delete-project", ProjectController.DeleteProject)
	//ProjectGroup.POST("/get-members", ProjectController.GetProjectMembers)
	//ProjectGroup.POST("/add-member", ProjectController.ProjectAddMember)
	//ProjectGroup.POST("/edit-member-power", ProjectController.ProjectEditMemberPower)
	//ProjectGroup.POST("/remove-member", ProjectController.ProjectRemoveMember)
	//
	//UserGroup.POST("/login", UserController.Login)
	//UserGroup.POST("/create-user", UserController.CreateUser)
	//UserGroup.POST("/search-user", UserController.SearchUserByUserName)
	//ImageGroup.GET("/getPlatFormProjectLists", UserController.GetProjects)
	//
	////镜像接口
	//ImageGroup.POST("/commonList", ImageController.GetCommonImages)
	//ImageGroup.POST("/buildImage", ImageController.BuildImage)
	//ImageGroup.POST("/imageVersionList", ImageController.GetArtifacts)
	//ImageGroup.GET("/pushImage", ImageController.PushImage)
	//ImageGroup.POST("/buildAndPushImage", ImageController.BuildAndPushImage)
	//ImageGroup.POST("/deleteImage", ImageController.DeleteImage)
	//ImageGroup.POST("/deleteRepo", ImageController.DeleteRepo)
	//
	//UserGroup.POST("/test", UserController.TestToken)

}
