package server

import "github.com/gin-gonic/gin"

func newRouter() *gin.Engine {
	g := gin.Default()
	g.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	return g
}
