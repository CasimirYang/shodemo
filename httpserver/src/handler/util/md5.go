package util

import (
	"crypto/md5"
	"fmt"
	"github.com/spf13/viper"
)

var md5Salt = viper.GetString("password.md5Salt")

func Md5Encode(password string) string {
	data := []byte(password + md5Salt)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}
