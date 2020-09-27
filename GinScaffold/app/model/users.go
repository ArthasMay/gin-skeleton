package model

import (
	"fmt"
	"go.uber.org/zap"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/yml_config"
	"log"
)

func CreateUserFactory(sqlType string) *userModel {
	if len(sqlType) == 0 {
		sqlType = yml_config.CreateYamlFactory().GetString("UseDbType")
	}

	dbDriver := CreateBaseSqlFactory(sqlType)
	if dbDriver != nil {
		return &userModel{BaseModel: dbDriver}
	}
	log.Fatal("userModel工厂初始化失败")
	return nil
}

type userModel struct {
	*BaseModel
	Id       int64  `json:"id"`
	UserName string `json:"username"`
	Pass     string `json:"-"`
	Phone    string `json:"phone"`
	RealName string `json:"realname"`
	Status   int    `json:"status"`
	Token    string `json:"-"`
}

// 用户注册
func (u *userModel) Register(username, pass, userIp string) bool {
	// 防止重复插入，但是 unique index 更加科学
	sql := "INSERT INTO tb_users(username, pass, last_login_ip) SELECT ?,?,? FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM tb_users WHERE username=?)"
	if u.ExecuteSql(sql, username, pass, userIp, username) > 0 {
		return true
	}
	return false
}

// 用户登录
func (u *userModel) Login(username, pass string) *userModel {
	sql := "SELECT id, username, pass, phone FROM tb_users WHERE username=?"
	err := u.QueryRow(sql, username).Scan(&u.Id, &u.UserName, &u.Pass, &u.Phone)
	if err == nil {
		// 验证密码
		fmt.Println(md5_encrypt.Base64Md5(pass))
		if len(u.Pass) > 0 && (u.Pass == md5_encrypt.Base64Md5(pass)) {
			return u
		}
	} else {
		variable.ZapLog.Error("根据账号查询单条记录出错:", zap.Error(err))
	}
	return nil
}

// 登录成功，记录 token
func (u *userModel) OauthLoginToken(userId int64, userToken string, expiresAt int64, clientIp string) bool {
	sql := "INSERT INTO `tb_oauth_access_tokens` (fr_user_id, `action_name`, token, expires_at, client_ip)" + "SELECT ?, 'login', ?, FROM_UNIXTIME(?), ? FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM `tb_oauth_access_tokens` a WHERE a.fr_user_id=? AND a.action_name='login' AND a.token=?)"
	//注意：token的精确度为秒，如果在一秒之内，一个账号多次调用接口生成的token其实是相同的，这样写入数据库，第二次的影响行数为0，知己实际上操作仍然是有效的。
	//所以这里的判断影响行数>=0 都是正确的，只有 -1 才是执行失败、错误
	if u.ExecuteSql(sql, userId, userToken, expiresAt, clientIp, userId, userToken) >= 0 {
		return true
	}
	return false
}

// 刷新token
func (u *userModel) OauthRefreshToken(userId, expiresAt int64, oldToken, newToken, clientIp string) bool {
	// sql := "UPDTAE `tb_oauth_access_tokens` SET token=?, expires_at=FROM_UNIXTIME(?), client_ip=?, update_at=NOW() WHERE fr_user_id=? AND token=?"
	// if u.ExecuteSql(sql, newToken, expiresAt, clientIp, userId, oldToken) > 0 {
	// 	return true
	// }
	// return false
	sql := "UPDATE tb_oauth_access_tokens SET token=?, expires_at=FROM_UNIXTIME(?), client_ip=?, updated_at=NOW() WHERE fr_user_id=? AND token=?"
	variable.ZapLog.Sugar().Info(sql, newToken, expiresAt, clientIp, userId, oldToken)
	if u.ExecuteSql(sql, newToken, expiresAt, clientIp, userId, oldToken) > 0 {
		return true
	}
	return false
}

// 校验token
func (u *userModel) OauthCheckTokenIsOk(userId int64, token string) bool {
	sql := "SELECT token FROM tb_oauth_access_tokens WHERE fr_user_id=? AND revoked=0 AND expires_at>NOW() ORDER BY update_at DESC LIMIT ?"
	rows := u.QuerySql(sql, userId, consts.JwtTokenOnlineUsers)
	
	if rows != nil {
		for rows.Next() {
			var tmpToken string
			err := rows.Scan(&tmpToken) 
			if err == nil {
				if tmpToken == token {
					_ = rows.Close()
					return true
				}
			}
		}
		_ = rows.Close()
	}
	return false
}
