package server

import (
	"github.com/PopescuStefanRadu/ent-demo/pkg"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server/controller"
	"github.com/gin-gonic/gin"
)

func NewRouter(app *pkg.App) *gin.Engine {
	g := gin.Default()

	userCtl := controller.User{UserService: app.Service}

	g.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	g.GET("/user/:id", userCtl.Get)
	g.POST("/user", userCtl.Create)
	g.PUT("/user/:id", userCtl.Update)
	g.DELETE("/user/:id", userCtl.Delete)
	g.POST("/search-users", userCtl.GetFiltered)

	return g
}
