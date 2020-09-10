package routers

import (
	"net/http"
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

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Api接口模块：Hello, World！")
	})

	//处理静态资源（不建议gin框架处理静态资源，参见 Public/readme.md 说明 ）
	router.Static("/public", "./public")             // 定义静态资源路由与实际目录映射关系
	router.StaticFile("/abcd", "./public/readme.md") // 可以根据文件名绑定需要返回的文件名
	
	vApi := router.Group("/api/v1/")
	{
	vApi := vApi.Group("home/")
		{
			vApi.GET("news", validatorFactory.Create(consts.ValidatorPrefix + "HomeNews"))
		}
	}
	return router
}
