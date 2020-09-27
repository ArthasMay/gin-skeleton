package web

import (
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/service/users/curd"
	userstoken "goskeleton/app/service/users/token"
	"goskeleton/app/utils/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Users struct {
}

// 注册: Controller的业务处理
func (u *Users) Register(context *gin.Context) {
	//  由于本项目骨架已经将表单验证器的字段(成员)绑定在上下文，因此可以按照 GetString()、GetInt64()、GetFloat64（）等快捷获取需要的数据类型，注意：相关键名都是小写
	// 当然也可以通过gin框架的上下文原始方法获取，例如： context.PostForm("name")
	name := context.GetString(consts.ValidatorPrefix + "name")
	pass := context.GetString(consts.ValidatorPrefix + "pass")
	userIp := context.ClientIP()

	if curd.CreateUserCurdFactory().Register(name, pass, userIp) {
		response.ReturnJson(
			context,
			http.StatusOK,
			consts.CurdStatusOkCode,
			consts.CurdStatusOkMsg,
			"",
		)
	} else {
		response.ReturnJson(
			context,
			http.StatusOK,
			consts.CurdCreatFailCode,
			consts.CurdCreatFailMsg,
			"",
		)
	}
}

// 登录: Controller 的业务处理
func (u *Users) Login(context *gin.Context) {
	// 1. 获取参数
	name := context.GetString(consts.ValidatorPrefix + "name")
	pass := context.GetString(consts.ValidatorPrefix + "pass")
	phone := context.GetString(consts.ValidatorPrefix + "phone")

	// 2. 查询 user 表，判断是否存在
	userModel := model.CreateUserFactory("").Login(name, pass)
	if userModel != nil {
		// 3. 若是查询到用户信息，生成token，并且把 token 落表 token 表中
		userTokenFactory := userstoken.CreateUserTokenFactory()
		if userToken, err := userTokenFactory.GenerateToken(userModel.Id, userModel.UserName, userModel.Phone, consts.JwtTokenCreatedExpireAt); err == nil {
			if userTokenFactory.RecordUserToken(userToken, context.ClientIP()) {
				data := gin.H{
					"userId":    userModel.Id,
					"name":      name,
					"realName":  userModel.RealName,
					"phone":     phone,
					"token":     userToken,
					"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
				}
				response.ReturnJson(
					context,
					http.StatusOK,
					consts.CurdStatusOkCode,
					consts.CurdStatusOkMsg,
					data,
				)
				return
			}
		}
	}
	response.ReturnJson(
		context,
		http.StatusOK,
		consts.CurdLoginFailCode,
		consts.CurdLoginFailMsg,
		"",
	)
}

// 刷新token
func (u *Users) RefreshToken(context *gin.Context) {
	token := context.GetString(consts.ValidatorPrefix + "token")
	if newToken, ok := userstoken.CreateUserTokenFactory().RefreshUserToken(token, context.ClientIP()); ok {
		res := gin.H{
			"token": newToken,
		}
		response.ReturnJson(
			context,
			http.StatusOK,
			consts.CurdStatusOkCode,
			consts.CurdStatusOkMsg,
			res,
		)
	} else {
		response.ReturnJson(
			context,
			http.StatusOK,
			consts.CurdRefreshTokenFailCode,
			consts.CurdRefreshTokenFailMsg,
			"",
		)
	}
}


