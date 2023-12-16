package server

import (
	"net/http"

	"github.com/PopescuStefanRadu/ent-demo/pkg/app"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server/controller"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(app *app.App) *gin.Engine {
	g := gin.Default()

	userCtl := controller.User{UserService: app.Service}
	errorHandler := &middleware.ErrorHandler{Logger: app.Logger}

	grp := g.Use(errorHandler.HandleErrors)

	grp.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	grp.GET("/user/:id", userCtl.Get)
	grp.POST("/user", userCtl.Create)
	grp.PUT("/user/:id", userCtl.Update)
	grp.DELETE("/user/:id", userCtl.Delete)
	grp.POST("/search-users", userCtl.GetFiltered)

	return g
}
