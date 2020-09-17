package web

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
)

type Users struct {
}

func (u *Users) Register(context *gin.Context) {
	//  由于本项目骨架已经将表单验证器的字段(成员)绑定在上下文，因此可以按照 GetString()、GetInt64()、GetFloat64（）等快捷获取需要的数据类型，注意：相关键名都是小写
	// 当然也可以通过gin框架的上下文原始方法获取，例如： context.PostForm("name") 
	name := context.GetString(consts.ValidatorPrefix + "name")
	pass := context.GetString(consts.ValidatorPrefix + "pass")
	userIp := context.ClientIP()
}