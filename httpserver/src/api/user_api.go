package api

import (
	"github.com/gin-gonic/gin"
	"httpserver/handler"
	"httpserver/handler/util"
)

func RouteUser(router *gin.Engine) {

	router.Use(util.JwtMiddleware())

	router.POST("/uc/login", handler.Login())
	router.GET("/uc/auth/getUser", handler.GetUser())
	router.POST("/uc/auth/editUser", handler.EditUser())
	router.POST("/uc/auth/uploadProfile", handler.UploadProfile())
}
