package users

import (
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/web"
	"net/http"
	"goskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

type Show struct {
	Base
	// 为啥 int64不行
	Page   int64 `form:"page" json:"page" binding:"required,gt=0"` // 必填，page > 0
	Limits int64 `form:"limits" json:"limits" binding:"required,gt=0"` // 必填，limits > 0
}

func (s Show) CheckParams(context *gin.Context) {
	// TODO: 这里值传递和引用传递的实现区别
	if err := context.ShouldBind(&s); err != nil {
		errs := gin.H{
			"tips": "UserShow参数校验失败，参数不符合规定，name（不能为空）、page的值(>0)、limits的值（>0)",
			"err": err.Error(),
		}
		response.ReturnJson(context, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, errs)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(s, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ReturnJson(
			context,
			http.StatusInternalServerError,
			consts.ServerOccurredErrorCode,
			consts.ServerOccurredErrorMsg,
			"",
		)
	} else {
		(&web.Users{}).Show(extraAddBindDataContext)
	}
}