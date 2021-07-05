package main

import (
	"github.com/gin-gonic/gin"
	"httpserver/api"
)

//var log = logrus.New()

func main() {
	//// 为当前logrus实例设置消息的输出，同样地，
	//// 可以设置logrus实例的输出到任意io.writer
	//log.Out = os.Stdout
	//
	//// 为当前logrus实例设置消息输出格式为json格式。
	//// 同样地，也可以单独为某个logrus实例设置日志级别和hook，这里不详细叙述。
	//log.Formatter = &logrus.JSONFormatter{}

	//log.Info("======= log =======")

	router := gin.Default()

	//router.Use(loggerToFile())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	api.RouteUser(router)

	router.Run(":8083")
}
