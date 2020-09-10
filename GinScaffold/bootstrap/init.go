package bootstrap

import (
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/validator/common/register_validator"
	"log"
	"os"
)

// 检查项目必须的非编译目录是否存在，避免编译后调用的时候缺失相关目录
func checkRequiredFolders() {
	// 1.检查配置文件是否存在
	mrp := os.Getenv("mrp")
	if _, err := os.Stat(variable.BasePath + mrp + "/config/config.yml"); err != nil {
		log.Fatal(my_errors.ErrorsConfigYamlNotExists + err.Error())
	}
	// 2.检查public目录是否存在
	if _, err := os.Stat(variable.BasePath + mrp + "/public/"); err != nil {
		log.Fatal(my_errors.ErrorsPublicNotExists + err.Error())
	}

	// 3.检查storage/logs目录是否存在
	if _, err := os.Stat(variable.BasePath + mrp + "/storage/logs/"); err != nil {
		log.Fatal(my_errors.ErrorsStorageLogsNotExists + err.Error())
	}
}

func init() {
	// 1. 初始化 项目根路径，参见 variable 常量包，相关路径：app\global\variable\variable.go

	// 2.检查配置文件以及日志目录等非编译性的必要条件
	checkRequiredFolders()

	// 3.初始化全局日志句柄，并载入日志钩子处理函数
	// variable.ZapLog = zap_factory.CreateZapFactory(sys_log_hook.ZapLogHandler)

	// 4.初始化表单验证器，将所有需要的验证器都注册到容器中
	register_validator.RegisterValidator()

	// 5.websocket Hub中心启动
}
