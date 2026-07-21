package router

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	_ "shafurui/docs"
	"shafurui/internal/config"
	"shafurui/internal/controller"
	"shafurui/internal/middleware"
	"shafurui/internal/pkg/telegram"
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
	// /api/swagger/index.html
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// /api/k4/index.html
	r.GET("/api/k4/*any", knife4goGin.WrapHandler(knife4goFiles.Handler))

	fmt.Printf("Go Version %v\n", runtime.Version())

	tgConfig := config.GetTelegramConfig()
	var telegramSender service.TelegramMessageSender
	if strings.TrimSpace(tgConfig.BotToken) != "" && strings.TrimSpace(tgConfig.ChatID) != "" {
		telegramSender = telegram.NewClient(tgConfig.BotToken, tgConfig.ChatID)
	}
	authService := service.NewAuthService(telegramSender)
	authController := controller.NewAuthController(authService)

	metaController := controller.NewMetaController()
	appRepo := persistence.NewAppRepository()
	appService := service.NewAppService(appRepo)
	appController := controller.NewAppController(appService)
	userRepo := persistence.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)
	videoRepo := persistence.NewVideoRepository()
	videoService := service.NewVideoService(videoRepo)
	videoController := controller.NewVideoController(videoService)

	v1 := r.Group("/api")

	v1.GET("/meta", metaController.GetMeta)

	appGroup := v1.Group("/app")
	{
		appGroup.POST("/getHelloInfo", appController.GetHelloInfo)
	}

	authGroup := v1.Group("/auth")
	authGroup.POST("/login", authController.Login)
	authGroup.POST("/refresh", authController.RefreshToken)

	userGroup := v1.Group("/user")
	userGroup.Use(middleware.JWTAuthMiddleware())
	userGroup.GET("/info", userController.GetCurrentUserInfo)

	videoGroup := v1.Group("/video")
	videoGroup.Use(middleware.JWTAuthMiddleware())
	videoGroup.GET("", videoController.ListVideos)
	v1.POST("/video/refresh", videoController.RefreshVideos)

	if videoDirPath := config.GetVideoDirPath(); strings.TrimSpace(videoDirPath) != "" {
		r.Static("/videos", videoDirPath)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})
	return r
}
