package curd

import (
	"goskeleton/app/model"
	"goskeleton/app/utils/md5_encrypt"
)

func CreateUserCurdFactory() *UserCurd {
	return &UserCurd{}
}

type UserCurd struct {
}

func (u *UserCurd) Register(name string, pass string, userIp string) bool {
	pass = md5_encrypt.Base64Md5(pass)
	return model.CreateUserFactory("").Register(name, pass, userIp)
}