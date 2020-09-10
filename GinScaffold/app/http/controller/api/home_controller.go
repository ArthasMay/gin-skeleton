package api

import (
	"goskeleton/app/global/consts"
	"goskeleton/app/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 1. 门户类首页新闻
type Home struct {
}

func (u *Home) News(context *gin.Context) {
	newsType := context.GetString(consts.ValidatorPrefix + "newsType")
	page := context.GetFloat64(consts.ValidatorPrefix + "page")
	limit := context.GetFloat64(consts.ValidatorPrefix + "limit")
	userIp := context.ClientIP()

	// 这里随便模拟一条数据返回
	fakeData := gin.H{
		"newsType": newsType,
		"page":     page,
		"limit":    limit,
		"userIp":   userIp,
		"title":    "门户首页公司新闻标题001",
		"content":  "门户新闻内容001",
	}
	response.ReturnJson(context, http.StatusOK, consts.CurdStatusOkCode, consts.CurdStatusOkMsg, fakeData)
}
