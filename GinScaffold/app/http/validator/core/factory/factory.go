package factory

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/core/container"
	"goskeleton/app/http/validator/core/interf"
	// "goskeleton/app/global/my_errors"
	// "goskeleton/app/global/variable"
	"log"
)

// 表单参数验证器工厂（勿修改）
func Create(key string) func(context *gin.Context) {
	log.Println("sss")
	if value := container.CreateContainersFactory().Get(key); value != nil {
		log.Println(value)
		if val, isOK := value.(interf.ValidatorInterface); isOK {
			return val.CheckParams
		}
	}
	
	// variable.ZapLog.Error(my_errors.ErrorsValidatorNotExists + ", 验证器模块：" + key)
	return nil
}