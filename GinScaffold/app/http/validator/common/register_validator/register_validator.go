package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/validator/api/home"
)

// 各个业务模块验证器都必须要进行注册（初始化），程序启动时会自动加载到容器
func RegisterValidator() {
	containers := container.CreateContainersFactory()

	// 注册门户表单参数验证器
	key := consts.ValidatorPrefix + "HomeNews"
	containers.Set(key, home.News{})
}