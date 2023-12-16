package controller

import (
	"net/http"

	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server/request"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server/response"
	"github.com/PopescuStefanRadu/ent-demo/pkg/user"
	"github.com/gin-gonic/gin"
)

type User struct {
	UserService *user.Service
}

func (ctl *User) Get(c *gin.Context) {
	q := struct {
		ID int `binding:"required" uri:"id"`
	}{}

	if err := c.ShouldBindUri(&q); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := ctl.UserService.GetUserByID(c, q.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Response[response.User]{Result: response.User(*res)})
}

func (ctl *User) Create(c *gin.Context) {
	var q request.CreateUser

	if err := c.ShouldBind(&q); err != nil {
		_ = c.Error(err)
		return
	}

	u := user.CreateUserParams(q)

	created, err := ctl.UserService.CreateUser(c, &u)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Response[response.User]{Result: response.User(*created)})
}

func (ctl *User) Update(c *gin.Context) {
	var (
		q request.UpdateUserURI
		b request.UpdateUserBody
	)

	if err := c.ShouldBindUri(&q); err != nil {
		_ = c.Error(err)
		return
	}

	if err := c.ShouldBind(&b); err != nil {
		_ = c.Error(err)
		return
	}

	u := user.UpdateUserParams{
		ID:       q.ID,
		Username: b.Username,
		Email:    b.Email,
	}

	updated, err := ctl.UserService.UpdateUser(c, &u)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Response[response.User]{Result: response.User(*updated)})
}

func (ctl *User) Delete(c *gin.Context) {
	q := struct {
		ID int `binding:"required" uri:"id"`
	}{}

	if err := c.ShouldBindUri(&q); err != nil {
		_ = c.Error(err)
		return
	}

	if err := ctl.UserService.DeleteUserByID(c, q.ID); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (ctl *User) GetFiltered(c *gin.Context) {
	var q request.GetFilteredUsers

	if err := c.ShouldBind(&q); err != nil {
		_ = c.Error(err)
		return
	}

	f := user.FindAllFilter(q)

	filtered, err := ctl.UserService.FindAllUsersByFilter(c, &f)
	if err != nil {
		_ = c.Error(err)
		return
	}

	mapped := make([]response.User, len(filtered))
	for i, u := range filtered {
		mapped[i] = response.User(u)
	}

	c.JSON(http.StatusOK, response.Response[[]response.User]{Result: mapped})
}
