package router

import (
	"goIM/middlewear"
	"goIM/service"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	ginEngine := gin.Default()

	// v1 组
	v1 := ginEngine.Group("v1")

	// 用户模块 后续有用户api 就放在user组下
	user := v1.Group("user")
	{
		user.POST("/register", service.UserService.Register)
		user.POST("/login", service.UserService.Login)
		user.POST("/list", middlewear.JWY(), service.UserService.List)
		user.POST("/update", middlewear.JWY(), service.UserService.Update)
		user.POST("/delete", middlewear.JWY(), service.UserService.Delete)

		user.GET("/chat", service.MessageService.Chat)
	}

	relation := v1.Group("relation").Use(middlewear.JWY())
	{
		relation.POST("/list", service.RelationService.FriendList)
		relation.POST("/add", service.RelationService.AddFriend)
	}

	group := v1.Group("group").Use(middlewear.JWY())
	{
		group.POST("/list", service.GroupService.GroupList)
		group.POST("/create", service.GroupService.CreateGroup)
		group.POST("/join", service.GroupService.JoinGroup)
	}

	asset := v1.Group("asset").Use(middlewear.JWY())
	{
		asset.POST("/upImage", service.AssetService.UpImage)
	}

	return ginEngine
}
