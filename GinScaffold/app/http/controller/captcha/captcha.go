package chaptcha

import (
	"bytes"
	"goskeleton/app/global/consts"
	"goskeleton/app/utils/response"
	"net/http"
	"path"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

type Captcha struct {
	Id      string `json:"id"`
	ImgUrl  string `json:"img_url"`
	Refresh string `json:"refresh"`
	Verify  string `json:"verify"`
}

// 生成验证码ID
func (c *Captcha) GenerateId(context *gin.Context) {
	// 设置验证码的数字长度 （个数）
	length := 4
	captchaId := captcha.NewLen(length)
	c.Id = captchaId
	c.ImgUrl = "/captcha/" + captchaId + ".png"
	c.Refresh = c.ImgUrl + "?reload=1"
	c.Verify = "/captcha/" + captchaId + "/这里替换为正确的验证码进行验证"
	response.ReturnJson(
		context,
		http.StatusOK,
		200,
		"验证码信息",
		c,
	)
}

// 获取验证码信息
func (c *Captcha) GetImg(context *gin.Context) {
	captchaId := context.Param("captchaId")
	_, file := path.Split(context.Request.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || captchaId == "" {
		response.ReturnJson(
			context,
			http.StatusBadRequest,
			consts.CaptchaGetParamsInvalidCode,
			consts.CaptchaGetParamsInvalidMsg,
			nil,
		)
		return
	}

	if context.Query("reload") != "" {
		captcha.Reload(id)
	}

	context.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	context.Header("Param", "no-cache")
	context.Header("Expires", "0")

	var vBytes bytes.Buffer
	if ext == ".png" {
		context.Header("Content-type", "image/png")
		_ = captcha.WriteImage(&vBytes, id, captcha.StdWidth, captcha.StdHeight)
		http.ServeContent(context.Writer, context.Request, id+ext, time.Time{}, bytes.NewReader(vBytes.Bytes()))
	}
}

// 校验验证码
func (c *Captcha) CheckCode(context *gin.Context) {
	captchaId := context.Param("captchaId")
	value := context.Param("value")
	if captchaId == "" || value == "" {
		response.ReturnJson(
			context,
			http.StatusBadRequest,
			consts.CaptchaCheckParamsInvalidCode,
			consts.CaptchaCheckParamsInvalidMsg,
			nil,
		)
	}

	if captcha.VerifyString(captchaId, value) {
		response.ReturnJson(
			context,
			http.StatusOK,
			consts.CaptchaCheckOkCode,
			consts.CaptchaCheckOkMsg,
			nil,
		)
	} else {
		response.ReturnJson(
			context,
			http.StatusBadRequest,
			consts.CaptchaCheckFailCode,
			consts.CaptchaCheckFailMsg,
			nil,
		)
	}
}