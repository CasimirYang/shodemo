package api

import (
	"github.com/gin-gonic/gin"
	"httpserver/handler"
	"httpserver/handler/util"
)

func RouteUser(router *gin.Engine) {

	router.Use(util.JwtMiddleware())

	router.POST("/uc/login", handler.Login())
	router.GET("/uc/getUser", handler.GetUser())
	router.POST("/uc/editUser", handler.EditUser())
	router.POST("/uc/uploadProfile", handler.UploadProfile())
}
