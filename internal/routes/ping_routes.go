package routes

import "github.com/gin-gonic/gin"

func PingRoutes(api *gin.RouterGroup) {
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
