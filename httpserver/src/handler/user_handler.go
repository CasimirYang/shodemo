package handler

import (
	"errors"
	"github.com/CasimirYang/share"
	"github.com/gin-gonic/gin"
	"httpserver/handler/rpc"
	"httpserver/handler/util"
	"httpserver/handler/vo"
	"io"
	"net/http"
	"os"
	"path"
)

func Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		var json vo.LoginRequestVO
		err := c.ShouldBindJSON(&json)
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
		err := c.ShouldBindJSON(&json)
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
		filePath, err := generateFile(c)
		if err != nil {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: share.InvalidParams, Message: err.Error()})
			return
		}
		userInfoReply, err := rpc.EditUser(c.GetInt64("uid"), nil, &filePath)
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
		share.SugarLogger.Error(err)
		return "", err
	}
	ext := path.Ext(fileHeader.Filename)
	if ext != ".png" {
		return "", errors.New("need png format")
	}
	share.SugarLogger.Infof("'%s' uploaded!", fileHeader.Filename)

	out, err := os.Create("./" + fileHeader.Filename)
	if err != nil {
		share.SugarLogger.Error(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		share.SugarLogger.Error(err)
		return "", err
	}
	return out.Name(), nil
}
