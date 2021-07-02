package main

import (
	"github.com/gin-gonic/gin"
	"httpserver/api"
)

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	api.RouteUser(router)

	router.Run(":8081")
}
