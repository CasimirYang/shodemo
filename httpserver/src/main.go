package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"httpserver/api"
)

func main() {
	router := gin.Default()
	router.Use(api.AccessLogHandler())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	api.RouteUser(router)

	router.Run(viper.GetString("port"))
}
