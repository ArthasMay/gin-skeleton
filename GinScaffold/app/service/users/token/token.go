package token

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"goskeleton/app/http/middleware/my_jwt"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/global/variable"
	"goskeleton/app/global/my_errors"
)

// userToken 小写的，但是通过返回值可以给外面引用
func CreateUserTokenFactory() *userToken {
	return &userToken {
		userJwt: my_jwt.CreateMyJWT(consts.JwtTokenSignKey),
	}
}

type userToken struct {
	userJwt *my_jwt.JwtSign
}

// 生成 token
func (u *userToken) GenerateToken(userid int64, username string, phone string, expiresAt int64) (token string, err error) {
	customClaims := my_jwt.CustomClaims{
		UserId: userid,
		Name: username,
		Phone: phone,
		// 特别注意，针对前文的匿名结构体，初始化的时候必须指定键名，并且不带 jwt. 否则报错：Mixture of field: value and value initializers
		StandardClaims: jwt.StandardClaims {
			NotBefore: time.Now().Unix() - 10,
			ExpiresAt: time.Now().Unix() + expiresAt,
		},
	}
	return u.userJwt.CreateToken(customClaims)
}

// 用户登录成功，记录token
func (u *userToken) RecordUserToken(userToken, clientIp string) bool {
	if customClaims, err := u.userJwt.ParseToken(userToken); err == nil {
		userid := customClaims.UserId
		expiresAt := customClaims.ExpiresAt
		return model.CreateUserFactory("").OauthLoginToken(userid, userToken, expiresAt, clientIp)
	} else {
		return false
	}
}

// 刷新token的有限期（默认 + 3600秒，参见常见配置）
func (u *userToken) RefreshUserToken(oldToken, clientIp string) (newToken string, isRefresh bool) {
	_, code := u.isNotExpired(oldToken)
	switch code {
	case consts.JwtTokenOK, consts.JwtTokenExpired:
		// 若token失效，更新token
		if newToken, err := u.userJwt.RefreshToken(oldToken, consts.JwtTokenRefreshExpireAt); err == nil {
			if customClaims, err := u.userJwt.ParseToken(newToken); err == nil {
				userId := customClaims.UserId
				expiresAt := customClaims.ExpiresAt
				if model.CreateUserFactory("").OauthRefreshToken(userId, expiresAt, oldToken, newToken, clientIp) {
					return newToken, true
				}
			}
		}
	case consts.JwtTokenInvalid:
		variable.ZapLog.Error(my_errors.ErrorsTokenInvalid)
	}
	return "", false
}

// 解析token，并且判断token是否失效
func (u *userToken) isNotExpired(token string) (*my_jwt.CustomClaims, int) {
	if customClaims, err := u.userJwt.ParseToken(token); err == nil {
		if time.Now().Unix() - customClaims.ExpiresAt < 0 {
			// 有效token: 未过期
			return customClaims, consts.JwtTokenOK
		} else {
			// 过期的token
			return customClaims, consts.JwtTokenExpired
		}
	}
	// 无效token
	return nil, consts.JwtTokenInvalid
}

// 校验 token 是否有效
func (u *userToken) IsEffective(token string) bool {
	customClaims, code := u.isNotExpired(token)
	if consts.JwtTokenOK == code {
		if model.CreateUserFactory("").OauthCheckTokenIsOk(customClaims.UserId, token) {
			return true
		}
	}
	return false
}
