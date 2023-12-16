package server

import (
	"errors"
	"github.com/PopescuStefanRadu/ent-demo/pkg"
	"github.com/PopescuStefanRadu/ent-demo/pkg/ent"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server/controller"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(app *pkg.App) *gin.Engine {
	g := gin.Default()

	userCtl := controller.User{UserService: app.Service}

	grp := g.Use(func(c *gin.Context) {
		c.Next()
		errs := c.Errors
		if errs == nil {
			return
		}

		// TODO do not use ent.NotFoundError, instead create a business error that wraps these cases.
		var e *ent.NotFoundError
		for _, err := range errs {
			if errors.As(err, &e) {
				c.JSON(http.StatusNotFound, response.Response[*struct{}]{
					Errors: map[string]response.Error{
						"global": {
							Code:    "NotFound",
							Message: "resource not found",
						},
					},
				})
			}
		}

		app.Logger.Error().Msg("TODO - Handle")
	})

	grp.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	grp.GET("/user/:id", userCtl.Get)
	grp.POST("/user", userCtl.Create)
	grp.PUT("/user/:id", userCtl.Update)
	grp.DELETE("/user/:id", userCtl.Delete)
	grp.POST("/search-users", userCtl.GetFiltered)

	return g
}
