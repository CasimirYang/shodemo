package main

import (
	"github.com/CasimirYang/share"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"httpserver/api"
	"os"
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

	err := router.Run(viper.GetString("port"))
	if err != nil {
		share.SugarLogger.Error(err)
		os.Exit(1)
	}
}
