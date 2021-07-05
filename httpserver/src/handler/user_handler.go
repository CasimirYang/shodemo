package handler

import (
	"fmt"
	"github.com/CasimirYang/share"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"httpserver/handler/rpc"
	"httpserver/handler/util"
	"httpserver/handler/vo"
	"io"
	"log"
	"net/http"
	"os"
)

func Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		span := opentracing.StartSpan("/uc/login") // Start a span using the global, in this case noop, tracer
		span.LogKV("event", "vvv")
		defer span.Finish()

		var json vo.LoginRequestVO
		err := c.BindJSON(&json)
		if err == nil {
			userInfoReply, err := rpc.Login(json.UserName, util.Md5Encode(json.Password))
			if err != nil {
				c.JSON(http.StatusOK, vo.CommonResponseVO{Code: share.SystemError})
				return
			}
			var response vo.CommonResponseVO
			if userInfoReply.Code == share.Success {
				userInfo := userInfoReply.UserInfo
				userInfoVO := vo.UserInfoVO{userInfo.GetUserName(),
					userInfo.GetNickName(),
					userInfo.GetPassword(),
					userInfo.GetProfile()}
				token, err := util.GenerateToken(userInfo.GetUid())
				if err != nil {
					c.JSON(http.StatusOK, vo.CommonResponseVO{Code: share.SystemError})
					return
				}
				response = vo.CommonResponseVO{share.Success, &vo.UserResponseVO{token, &userInfoVO}}
			} else {
				response = vo.CommonResponseVO{Code: int(userInfoReply.Code)}
			}
			c.JSON(http.StatusOK, response)
		} else {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: share.InvalidParams, Message: err.Error()})
		}
	}
}

func GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		userInfoReply, err := rpc.GetUser(c.GetInt64("uid"))
		if err != nil {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: share.SystemError})
			return
		}
		var response vo.CommonResponseVO
		if userInfoReply.Code == share.Success {
			userInfo := userInfoReply.UserInfo
			userInfoVO := vo.UserInfoVO{userInfo.GetUserName(),
				userInfo.GetNickName(),
				userInfo.GetPassword(),
				userInfo.GetProfile()}
			response = vo.CommonResponseVO{share.Success, &vo.UserResponseVO{UserInfo: &userInfoVO}}
		} else {
			response = vo.CommonResponseVO{Code: int(userInfoReply.Code)}
		}
		c.JSON(http.StatusOK, response)
	}
}

func EditUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var json vo.UpdateRequestVO
		err := c.BindJSON(&json)
		if err == nil {
			userInfoReply, err := rpc.EditUser(c.GetInt64("uid"), &json.NickName, nil)
			if err != nil {
				c.JSON(http.StatusOK, vo.CommonResponseVO{Code: share.SystemError})
				return
			}
			var response vo.CommonResponseVO
			if userInfoReply.Code == share.Success {
				response = vo.CommonResponseVO{Code: share.Success}
			} else {
				response = vo.CommonResponseVO{Code: int(userInfoReply.Code)}
			}
			c.JSON(http.StatusOK, response)
		} else {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: share.InvalidParams, Message: err.Error()})
		}
	}
}

func UploadProfile() func(c *gin.Context) {
	return func(c *gin.Context) {
		path, err := generateFile(c)
		if err != nil {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: share.InvalidParams, Message: err.Error()})
			return
		}
		userInfoReply, err := rpc.EditUser(c.GetInt64("uid"), nil, &path)
		if err != nil {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: share.SystemError})
			return
		}
		var response vo.CommonResponseVO
		if userInfoReply.Code == share.Success {
			response = vo.CommonResponseVO{Code: share.Success}
		} else {
			response = vo.CommonResponseVO{Code: int(userInfoReply.Code)}
		}
		c.JSON(http.StatusOK, response)
	}
}

func generateFile(c *gin.Context) (string, error) {
	//limit 2m
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 2*1024*1024)
	// 拿到这个文件
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Printf("'%s' uploaded!", fileHeader.Filename)

	out, err := os.Create("./" + fileHeader.Filename + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return out.Name(), nil
}
