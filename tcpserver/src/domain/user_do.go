package domain

import (
	"database/sql"
	"errors"
	commonLog "git.garena.com/jinghua.yang/entry-task-common/log"
	"tcpserver/infrastructure/mysql"
	"tcpserver/infrastructure/po"
	"tcpserver/infrastructure/redis"
)

type UserDO struct {
	Id       int64
	UserName string
	NickName string
	Password string
	Profile  string
}

var ErrSystem = errors.New("system error")
var ErrNoData = errors.New("could not find match data")

//Login login and cache user
func Login(userName, password string) (*UserDO, error) {
	userPO, err := mysql.GetUser(userName, password)
	if err == sql.ErrNoRows {
		return nil, ErrNoData
	} else if err != nil {
		commonLog.SugarLogger.Error(err)
		return nil, ErrSystem
	}
	redis.CacheUser(userPO.Id, *userPO)
	userDO := UserDO(*userPO)
	return &userDO, nil
}

//GetUser getUser and delete cache
func GetUser(uid int64) (*UserDO, error) {
	var userPO *po.UserPO
	var err error
	userPO, err = redis.GetUser(uid)
	if err != nil {
		commonLog.SugarLogger.Error(err)
		return nil, ErrSystem
	} else if userPO == nil {
		userPO, err = mysql.GetUserByUid(uid)
		if err == sql.ErrNoRows {
			return nil, ErrNoData
		} else if err != nil {
			return nil, ErrSystem
		}
		redis.CacheUser(userPO.Id, *userPO)
	}
	userDO := UserDO(*userPO)
	return &userDO, nil
}

//UpdateUserByUid updateUser and delete cache
func UpdateUserByUid(uid int64, nickName, profile string) error {
	var err error
	if len(nickName) != 0 {
		err = mysql.UpdateNickName(uid, nickName)
	} else if len(profile) != 0 {
		err = mysql.UpdateProfile(uid, profile)
	}
	if err != nil {
		commonLog.SugarLogger.Error(err)
		return ErrSystem
	}
	redis.DeleteCache(uid)
	return nil
}
