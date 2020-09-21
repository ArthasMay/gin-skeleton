package model

import (
	// "go.uber.org/zap"
	// "goskeleton/app/global/consts"
	// "goskeleton/app/global/variable"
	// "goskeleton/app/utils/md5_encrypt"
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

// 简单的注册
func (u *userModel) Register(username, pass, userIp string) bool {
	// 防止重复插入，但是 unique index 更加科学
	sql := "INSERT INTO tb_users(username, pass, last_login_ip) SELECT ?,?,? FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM tb_users WHERE username=?)"
	if u.ExecuteSql(sql, username, pass, userIp, username) > 0 {
		return true
	}
	return false
}
