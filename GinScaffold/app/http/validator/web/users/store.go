package users

import (
	"goskeleton/app/http/controller/web"
	"goskeleton/app/http/validator/core/data_transfer"
	"net/http"
	"goskeleton/app/global/consts"
	"goskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

type Store struct {
	Base
	Pass     string `form:"pass" json:"pass" binding:"required,min=6"`
	RealName string `form:"real_name" json:"real_name" binding:"required,min=2"`
	Phone    string `form:"phone" json:"phone" binding:"required,len=11"`
	Remark   string `form:"remark" json:"remark"`
}

func (s Store) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&s); err != nil {
		errs := gin.H {
			"tips": "UserStore参数校验失败，参数校验失败，请检查name(>=1)、pass(>=6)、real_name(>=2)、phone(=11)",
			"err": err.Error(),
		}
		response.ReturnJson(context, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, errs)
		return
	}
	
	extraAddBindContext := data_transfer.DataAddContext(s, consts.ValidatorPrefix, context)
	if extraAddBindContext == nil {
		response.ReturnJson(
			context,
			http.StatusInternalServerError,
			consts.ServerOccurredErrorCode,
			consts.ServerOccurredErrorMsg,
			"",
		)
	} else {
		(&web.Users{}).Store(extraAddBindContext)
	}
}