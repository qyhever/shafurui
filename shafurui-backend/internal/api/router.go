package router

import (
	"fmt"
	"net/http"
	"runtime"

	_ "shafurui/docs"
	"shafurui/internal/config"
	"shafurui/internal/controller"
	"shafurui/internal/middleware"
	"shafurui/internal/repository/persistence"
	"shafurui/internal/service"

	"github.com/gin-gonic/gin"
	knife4goFiles "github.com/go-webtools/knife4go"
	knife4goGin "github.com/go-webtools/knife4go/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	isProd := config.IsProduction()
	// Gin 开启生产模式(默认是debug模式，会输出大量调试日志)
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 静态文件服务
	r.Static("/public", "./public")
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/api/k4/*any", knife4goGin.WrapHandler(knife4goFiles.Handler))

	fmt.Printf("Go Version %v\n", runtime.Version())

	authService := service.NewAuthService()
	authController := controller.NewAuthController(authService)

	metaController := controller.NewMetaController()
	appRepo := persistence.NewAppRepository()
	appService := service.NewAppService(appRepo)
	appController := controller.NewAppController(appService)

	v1 := r.Group("/api")

	v1.GET("/meta", metaController.GetMeta)

	appGroup := v1.Group("/app")
	authGroup := v1.Group("/auth")
	adminGroup := v1.Group("/admin")

	{
		appGroup.POST("/getHelloInfo", appController.GetHelloInfo)
	}

	authGroup.POST("/login", authController.Login)
	authGroup.POST("/refresh", authController.RefreshToken)
	adminProtectedGroup := adminGroup.Group("")
	adminProtectedGroup.Use(middleware.JWTAuthMiddleware())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})
	return r
}
