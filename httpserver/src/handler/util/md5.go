package util

import (
	"crypto/md5"
	"fmt"
)

var md5Salt = "sp001"

func Md5Encode(password string) string {
	data := []byte(password + md5Salt)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}
