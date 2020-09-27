package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/validator/api/home"
	"goskeleton/app/http/validator/web/users"
)

// 各个业务模块验证器都必须要进行注册（初始化），程序启动时会自动加载到容器
func RegisterValidator() {
	// 创建容器
	containers := container.CreateContainersFactory()

	// key 按照前缀 + 模块 + 验证动作 格式，将各个模块验证注册在容器
	var key string
	// Web Module
	key = consts.ValidatorPrefix + "UsersRegister"
	containers.Set(key, users.Register{})

	key = consts.ValidatorPrefix + "UserLogin"
	containers.Set(key, users.Login{})

	key = consts.ValidatorPrefix + "RefreshToken"
	containers.Set(key, users.RefreshToken{})

	// Api Module
	// 注册门户表单参数验证器
	key = consts.ValidatorPrefix + "HomeNews"
	containers.Set(key, home.News{})
}