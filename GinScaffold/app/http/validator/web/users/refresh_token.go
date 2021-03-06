package users

import (
	"strings"
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/web"
	"net/http"
)

type RefreshToken struct {
	Authorization string `json:"token" header:"Authorization" binding:"required,min=20"`
}

func (r RefreshToken) CheckParams(context *gin.Context) {
	if err := context.ShouldBindHeader(&r); err != nil {
		errs := gin.H{
			"tips": "Token参数校验失败，参数不符合规定，token 长度>=20",
			"err": err.Error(),
		}
		response.ReturnJson(
			context,
			http.StatusBadRequest,
			consts.ValidatorParamsCheckFailCode,
			consts.ValidatorParamsCheckFailMsg,
			errs,
		)
		return
	}
	
	token := strings.Split(r.Authorization, " ")
	if len(token) == 2 {
		context.Set(consts.ValidatorPrefix + "token", token[1])
		(&web.Users{}).RefreshToken(context)
	} else {
		errs := gin.H{
			"tips": "Token不合法，token请放置在header头部分，按照按=>键提交，例如：Authorization：Bearer 你的实际token....",
		}
		response.ReturnJson(context, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, errs)
	}
}