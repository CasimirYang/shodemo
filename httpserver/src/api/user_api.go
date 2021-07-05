package api

import (
	"github.com/gin-gonic/gin"
	"httpserver/handler"
)

func RouteUser(router *gin.Engine) {

	orderGroup := router.Group("/uc/auth")
	orderGroup.Use(handler.JwtMiddleware())

	router.POST("/uc/login", handler.Login())
	router.GET("/uc/auth/getUser", handler.GetUser())
	router.POST("/uc/auth/editUser", handler.EditUser())
	router.POST("/uc/auth/uploadProfile", handler.UploadProfile())
}
