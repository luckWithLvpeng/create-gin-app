package routers

import (
	"eme/controllers"
	"eme/middleware/auth"
	"eme/pkg/config"

	// docs is generated by Swag CLI, you have to import it.
	_ "eme/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func init() {
	InitAll()
}

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	if config.RunMode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/v1")
	apiv1.POST("/user/login", controllers.UserLogin)
	apiv1.POST("/user/refreshToken", controllers.UserRefreshToken)
	apiv1.Use(auth.Auth())
	{
		apiv1.GET("/user", controllers.UserGet)
		apiv1.GET("/user/:id", controllers.UserGetByID)
		apiv1.DELETE("/user/:id", controllers.UserDeleteByID)
		apiv1.PUT("/user/:id", controllers.UserUpdate)
		apiv1.POST("/user/logout", controllers.UserLogout)
		apiv1.POST("/user", controllers.UserAdd)

		apiv1.GET("/role", controllers.RoleGet)
		apiv1.GET("/role/:id", controllers.RoleGetByID)
		apiv1.POST("/role", controllers.RoleAdd)
		apiv1.PUT("/role/:id", controllers.RoleUpdate)
		apiv1.DELETE("/role/:id", controllers.RoleDel)
	}
	return r
}
