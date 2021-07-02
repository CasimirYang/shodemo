package handler

import (
	"fmt"
	"github.com/CasimirYang/share"
	"github.com/gin-gonic/gin"
	"httpserver/handler/rpc"
	"io"
	"log"
	"net/http"
	"os"
)

type LoginRequestJson struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateRequestJson struct {
	NickName string `json:"nickName" binding:"required"`
}

type MessageVO struct {
	Token    string      `json:"token,omitempty"`
	UserInfo *UserInfoVO `json:"userInfo"`
}

type UserInfoVO struct {
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
	Password string `json:"password"`
	Profile  string `json:"profile"`
}

func Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		var json LoginRequestJson
		if c.BindJSON(&json) == nil {
			userInfoReply, err := rpc.Login(json.UserName, json.Password)
			if err != nil {
				c.JSON(http.StatusOK, ResponseVO{Code: share.SystemError})
				return
			}
			var response ResponseVO
			if userInfoReply.Code == share.Success {
				userInfo := userInfoReply.UserInfo
				userInfoVO := UserInfoVO{userInfo.GetUserName(),
					userInfo.GetNickName(),
					userInfo.GetPassword(),
					userInfo.GetProfile()}
				token, err := generateToken(userInfo.GetUid())
				if err != nil {
					c.JSON(http.StatusOK, ResponseVO{Code: share.SystemError})
					return
				}
				response = ResponseVO{share.Success, &MessageVO{token, &userInfoVO}}
			} else {
				response = ResponseVO{Code: int(userInfoReply.Code)}
			}
			c.JSON(http.StatusOK, response)
		}
	}
}

func GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		userInfoReply, err := rpc.GetUser(c.GetInt64("uid"))
		if err != nil {
			c.JSON(http.StatusOK, ResponseVO{Code: share.SystemError})
			return
		}
		var response ResponseVO
		if userInfoReply.Code == share.Success {
			userInfo := userInfoReply.UserInfo
			userInfoVO := UserInfoVO{userInfo.GetUserName(),
				userInfo.GetNickName(),
				userInfo.GetPassword(),
				userInfo.GetProfile()}
			response = ResponseVO{share.Success, &MessageVO{UserInfo: &userInfoVO}}
		} else {
			response = ResponseVO{Code: int(userInfoReply.Code)}
		}
		c.JSON(http.StatusOK, response)
	}
}

func EditUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var json UpdateRequestJson
		if c.BindJSON(&json) == nil {
			userInfoReply, err := rpc.EditUser(c.GetInt64("uid"), &json.NickName, nil)
			if err != nil {
				c.JSON(http.StatusOK, ResponseVO{Code: share.SystemError})
				return
			}
			var response ResponseVO
			if userInfoReply.Code == share.Success {
				response = ResponseVO{Code: share.Success}
			} else {
				response = ResponseVO{Code: int(userInfoReply.Code)}
			}
			c.JSON(http.StatusOK, response)
		}
	}
}

func UploadProfile() func(c *gin.Context) {
	return func(c *gin.Context) {
		path := generateFile(c)
		userInfoReply, err := rpc.EditUser(c.GetInt64("uid"), nil, &path)
		if err != nil {
			c.JSON(http.StatusOK, ResponseVO{Code: share.SystemError})
			return
		}
		var response ResponseVO
		if userInfoReply.Code == share.Success {
			response = ResponseVO{Code: share.Success}
		} else {
			response = ResponseVO{Code: int(userInfoReply.Code)}
		}
		c.JSON(http.StatusOK, response)
	}
}

func generateFile(c *gin.Context) string {
	//todo 文件大小限制
	// 拿到这个文件
	file, fileHeader, _ := c.Request.FormFile("file")
	if file == nil {
		//todo no support
	}
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", fileHeader.Filename))

	out, err := os.Create("./" + fileHeader.Filename + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	return out.Name()
}
