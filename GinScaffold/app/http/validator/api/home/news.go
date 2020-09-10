package home

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/api"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type News struct {
	// 验证规则：必填，最小长度为1
	NewType string `form:"newsType" json:"newsType" binding:"required,min=1"`
	// 验证规则：必填，最小值为1(float类型，min=1代表最小值为1)
	Page 	float64 `form:"page" json:"page" binding:"required,min=1"`
	// 验证规则：必填，最小值为1
	Limit  	float64 `form:"limit" json:"limit" binding:"required,min=1"`
}

func (n News) CheckParams(context *gin.Context) {
	// 1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		errs := gin.H{
			"tips": "HomeNews参数校验失败，参数不符合规定，newsType(长度>=1)、page>=1、limit>=1,请按照规则自己检查",
			"err":  err.Error(),
		}
		response.ReturnJson(context, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, errs)
		return
	}

	// 该函数主要是将绑定的数据以 键=>值对形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ReturnJson(context, http.StatusInternalServerError, consts.ServerOccurredErrorCode, consts.ServerOccurredErrorMsg+",HomeNews表单验证器json化失败", "")
	} else {
		(&api.Home{}).News(extraAddBindDataContext)
	}
}
