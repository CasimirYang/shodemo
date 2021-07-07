package main

import (
	_ "git.garena.com/jinghua.yang/entry-task-common/config"
	commonLog "git.garena.com/jinghua.yang/entry-task-common/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"httpserver/api"
	"net/http"
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

	router.StaticFS("/static", http.Dir("./static"))

	api.RouteUser(router)

	err := router.Run(viper.GetString("port"))
	if err != nil {
		commonLog.SugarLogger.Error(err)
		os.Exit(1)
	}
}
