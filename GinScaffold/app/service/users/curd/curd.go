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

func (u *UserCurd) Register(name string, pass string, userIp string) {
	pass = md5_encrypt.Base64Md5(pass)
	
}