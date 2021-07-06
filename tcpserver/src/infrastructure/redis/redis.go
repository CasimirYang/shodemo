package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"strconv"
	"tcpserver/infrastructure/po"
	"time"
)

var ctx = context.Background()

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.url"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

const UserCacheKey = "account_user_id_"

func CacheUser(uid int64, user po.UserPO) {
	idKey := UserCacheKey + strconv.FormatInt(uid, 10)
	rdb.HSet(ctx, idKey, "id", user.Id)
	rdb.HSet(ctx, idKey, "userName", user.UserName)
	rdb.HSet(ctx, idKey, "nickName", user.NickName)
	rdb.HSet(ctx, idKey, "password", user.Password)
	rdb.HSet(ctx, idKey, "profile", user.Profile)
	rdb.Expire(ctx, idKey, time.Hour)
}

func GetUser(uid int64) (*po.UserPO, error) {
	idKey := UserCacheKey + strconv.FormatInt(uid, 10)
	val, err := rdb.HGetAll(ctx, idKey).Result()
	if err != nil {
		return nil, err
	}
	if len(val) == 0 {
		return nil, nil
	}
	user := po.UserPO{Id: uid, UserName: val["userName"], NickName: val["nickName"], Password: val["password"], Profile: val["profile"]}
	return &user, nil
}

func DeleteCache(uid int64) {
	idKey := UserCacheKey + strconv.FormatInt(uid, 10)
	rdb.Del(ctx, idKey)
}
