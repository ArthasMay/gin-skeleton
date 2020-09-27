package authorization

import (
	"strings"
	"net/http"
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/global/my_errors"
	userstoken "goskeleton/app/service/users/token"
	"goskeleton/app/utils/response"
)

type HeaderParams struct {
	Authorization string `header:"Authorization"`
}

func CheckAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		headerParams := HeaderParams{}
		
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			variable.ZapLog.Error(
				my_errors.ErrorsValidatorBindParamsFail,
				zap.Error(err),
			)
			context.Abort()
		}

		if len(headerParams.Authorization) >= 20 {
			token := strings.Split(headerParams.Authorization, " ")
			if len(token) == 2 && len(token[1]) >= 20 {
				tokenIsEffective := userstoken.CreateUserTokenFactory().IsEffective(token[1])
				if tokenIsEffective {
					context.Next()
				} else {
					response.ReturnJson(
						context,
						http.StatusUnauthorized,
						http.StatusUnauthorized,
						my_errors.ErrorsNoAuthorization,
						"",
					)
					context.Abort()
				}
			}
		} else {
			response.ReturnJson(
				context, 
				http.StatusUnauthorized, 
				http.StatusUnauthorized, 
				my_errors.ErrorsNoAuthorization, 
				"",
			)
			//暂停执行
			context.Abort()
		}
	}
}