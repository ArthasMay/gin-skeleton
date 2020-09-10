package data_transfer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/validator/core/interf"
)

// 将验证器成员（字段）绑定到数据传输上下文，方便控制器获取
/**
 * validatorInterface 实现验证器接口的结构体
 * 
 */
func DataAddContext(
	validatorInterface interf.ValidatorInterface,
	extraAddDataPrefix string,
	context *gin.Context) *gin.Context {
	var tmpJson interface{}

	if tmpBytes, err1 := json.Marshal(validatorInterface); err1 == nil {
		if err2 := json.Unmarshal(tmpBytes, &tmpJson); err2 == nil {
			if value, ok := tmpJson.(map[string]interface{}); ok {
				for k, v := range value {
					context.Set(extraAddDataPrefix+k, v)
				}
				return context
			}
		}
	}
	return nil
}