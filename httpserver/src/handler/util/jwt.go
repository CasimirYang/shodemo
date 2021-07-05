package util

import (
	"github.com/CasimirYang/share"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"httpserver/handler/vo"

	"net/http"
	"time"
)

var jwtSecret = []byte("sho001")

type Claims struct {
	Uid int64 `json:"uid"`
	jwt.StandardClaims
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = share.Success
		token := c.GetHeader("token")
		var claims *Claims
		if token == "" {
			code = share.InvalidToken
		} else {
			var err error
			claims, err = parseToken(token)
			if err != nil {
				code = share.InvalidToken
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = share.InvalidToken
			}
		}
		if code != share.Success {
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
