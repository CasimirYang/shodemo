package handler

import (
	"errors"
	commonCode "git.garena.com/jinghua.yang/entry-task-common/code"
	commonLog "git.garena.com/jinghua.yang/entry-task-common/log"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
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
		var requestJson vo.LoginRequestVO
		err := c.ShouldBindJSON(&requestJson)
		if err != nil {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: commonCode.InvalidParams, Message: err.Error()})
			return
		}
		userInfoReply, err := rpc.Login(requestJson.UserName, util.Md5Encode(requestJson.Password))
		if err != nil {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: commonCode.SystemError})
			return
		}
		var response vo.CommonResponseVO
		if userInfoReply.Code == commonCode.Success {
			userInfo := userInfoReply.UserInfo
			userInfoVO := &vo.UserInfoVO{UserName: userInfo.GetUserName(),
				NickName: userInfo.GetNickName(),
				Profile:  userInfo.GetProfile()}
			token, err := util.GenerateToken(userInfo.GetUid())
			if err != nil {
				c.JSON(http.StatusOK, vo.CommonResponseVO{Code: commonCode.SystemError})
				return
			}
			response = vo.CommonResponseVO{Code: commonCode.Success, Message: &vo.UserResponseVO{Token: token, UserInfo: userInfoVO}}
		} else {
			response = vo.CommonResponseVO{Code: int(userInfoReply.Code)}
		}
		c.JSON(http.StatusOK, response)
	}
}

func GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		userInfoReply, err := rpc.GetUser(c.GetInt64("uid"))
		if err != nil {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: commonCode.SystemError})
			return
		}
		var response vo.CommonResponseVO
		if userInfoReply.Code == commonCode.Success {
			userInfo := userInfoReply.UserInfo
			userInfoVO := &vo.UserInfoVO{UserName: userInfo.GetUserName(),
				NickName: userInfo.GetNickName(),
				Profile:  userInfo.GetProfile()}
			response = vo.CommonResponseVO{Code: commonCode.Success, Message: &vo.UserResponseVO{UserInfo: userInfoVO}}
		} else {
			response = vo.CommonResponseVO{Code: int(userInfoReply.Code)}
		}
		c.JSON(http.StatusOK, response)
	}
}

func EditUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var requestJson vo.UpdateRequestVO
		err := c.ShouldBindJSON(&requestJson)
		if err != nil {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: commonCode.InvalidParams, Message: err.Error()})
		}
		userInfoReply, err := rpc.EditUser(c.GetInt64("uid"), &requestJson.NickName, nil)
		if err != nil {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: commonCode.SystemError})
			return
		}
		var response vo.CommonResponseVO
		if userInfoReply.Code == commonCode.Success {
			response = vo.CommonResponseVO{Code: commonCode.Success}
		} else {
			response = vo.CommonResponseVO{Code: int(userInfoReply.Code)}
		}
		c.JSON(http.StatusOK, response)
	}
}

func UploadProfile() func(c *gin.Context) {
	return func(c *gin.Context) {
		filePath, err := generateFile(c)
		if err != nil {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: commonCode.InvalidParams, Message: err.Error()})
			return
		}
		userInfoReply, err := rpc.EditUser(c.GetInt64("uid"), nil, &filePath)
		if err != nil {
			c.JSON(http.StatusOK, vo.CommonResponseVO{Code: commonCode.SystemError})
			return
		}
		var response vo.CommonResponseVO
		if userInfoReply.Code == commonCode.Success {
			response = vo.CommonResponseVO{Code: commonCode.Success}
		} else {
			response = vo.CommonResponseVO{Code: int(userInfoReply.Code)}
		}
		c.JSON(http.StatusOK, response)
	}
}

func generateFile(c *gin.Context) (string, error) {
	//limit 2m
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 2*1024*1024)
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		commonLog.SugarLogger.Error(err)
		return "", err
	}
	ext := path.Ext(fileHeader.Filename)
	if ext != ".png" {
		return "", errors.New("need png format")
	}
	newFileName := uuid.NewV4().String() + ".png"
	out, err := os.Create("static/" + newFileName)
	if err != nil {
		commonLog.SugarLogger.Error(err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			commonLog.SugarLogger.Error(err)
		}
	}(out)
	_, err = io.Copy(out, file)
	if err != nil {
		commonLog.SugarLogger.Error(err)
		return "", err
	}

	commonLog.SugarLogger.Infof("'%s' uploaded!", newFileName)

	return "http://localhost:8083/" + out.Name(), nil
}
