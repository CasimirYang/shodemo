package util

import (
	commonCode "git.garena.com/jinghua.yang/entry-task-common/code"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"httpserver/handler/vo"

	"net/http"
	"time"
)

var jwtSecret = []byte(viper.GetString("jwtSecret"))

type Claims struct {
	Uid int64 `json:"uid"`
	jwt.StandardClaims
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.RequestURI == "/uc/login" {
			c.Next()
			return
		}
		var code int
		code = commonCode.Success
		token := c.GetHeader("token")
		var claims *Claims
		if token == "" {
			code = commonCode.InvalidToken
		} else {
			var err error
			claims, err = parseToken(token)
			if err != nil {
				code = commonCode.InvalidToken
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = commonCode.InvalidToken
			}
		}
		if code != commonCode.Success {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: code})
			c.Abort()
			return
		}

		c.Set("uid", claims.Uid)
		c.Next()
	}
}

func GenerateToken(uid int64) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func parseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
