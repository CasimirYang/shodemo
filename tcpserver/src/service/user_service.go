package service

import (
	"tcpserver/domain"
)

func Login(userName, password string) (*domain.UserDO, error) {
	userDO, err := domain.Login(userName, password)
	return userDO, err
}

func GetUser(uid int64) (*domain.UserDO, error) {
	userDO, err := domain.GetUser(uid)
	return userDO, err
}

func UpdateUserByUid(uid int64, nickName, profile string) error {
	err := domain.UpdateUserByUid(uid, nickName, profile)
	return err
}
