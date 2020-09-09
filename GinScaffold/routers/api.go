package routers

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/pprof"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/yml_config"
	"goskeleton/app/http/middleware/cors"
	validatorFactory "goskeleton/app/http/validator/core/factory"
)

func InitApiRouter() *gin.Engine {
	var router *gin.Engine
	// 非调试模式（生产模式）日志写到日志文件
	if yml_config.CreateYamlFactory().GetBool("AppDebug") == false {
		// 1. 将日志写入日志文件
		gin.DisableConsoleColor()
		f, _ := os.Create(variable.BasePath + yml_config.CreateYamlFactory().GetString("Logs.GinLogName"))
		gin.DefaultWriter = io.MultiWriter(f)
		// 2. 如果有nginx前置做代理，基本不需要gin框架记录访问日志，开启下一行代码，屏蔽上面三行代码，性能提升 5%
		// gin.SetMode(gin.ReleaseMode)
		router = gin.Default()
	} else {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}

	// 根据配置文件设置跨域
	if yml_config.CreateYamlFactory().GetBool("HttpServer.AllowCrossDomain") {
		router.Use(cors.Next())
	}

	router.Static("/public", "./public")
	router.StaticFile("abcd", "./public/readme.md")
	
	vApi := router.Group("/api/v1/")
	{
		vApi := vApi.Group("home/")
		{
			vApi.GET("news", validatorFactory.Create(consts.ValidatorPrefix + "HomeNews"))
		}
	}
	return router
}
