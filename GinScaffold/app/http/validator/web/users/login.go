package users

import (
	"goskeleton/app/global/consts"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
	"goskeleton/app/http/controller/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Login struct {
	Base
	// 必填，密码长度范围：【3,20】闭区间
	Pass string `form:"pass" json:"pass" binding:"required,min=3,max=20"`
}

func (l Login) CheckParams(context *gin.Context) {
	// 1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&l); err != nil {
		errs := gin.H{
			"tips": "UserLogin参数校验失败，name长度(>=1)、pass长度(6,20)有问题，不允许登录",
			"err":  err.Error(),
		}
		response.ReturnJson(
			context,
			http.StatusBadRequest,
			consts.ValidatorParamsCheckFailCode,
			consts.ValidatorParamsCheckFailMsg,
			errs,
		)
	}

	// 该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(l, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ReturnJson(
			context,
			http.StatusInternalServerError,
			consts.ServerOccurredErrorCode,
			consts.ServerOccurredErrorMsg + "UserLogin表单验证器json化失败",
			"",
		)
	} else {
		// 验证完成，调用控制器，并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&web.Users{}).Login(extraAddBindDataContext)
	}
}
