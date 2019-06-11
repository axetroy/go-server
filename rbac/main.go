// Copyright 2019 Axetroy. All rights reserved. MIT license.
package rbac

import (
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/module/role/role_model"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/rbac/accession"
	"github.com/axetroy/go-server/rbac/role"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Controller struct {
	Roles []*role.Role
}

func New(uid string) (c *Controller, err error) {
	c = &Controller{}

	userInfo := user_model.User{
		Id: uid,
	}

	if err = database.Db.First(&userInfo).Error; err != nil {
		return
	}

	if len(userInfo.Role) == 0 {
		err = common_error.ErrNoPermission
		return
	}

	for _, roleName := range userInfo.Role {
		roleInfo := role_model.Role{
			Name: roleName,
		}

		if err = database.Db.First(&roleInfo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return
		}

		r := role.New(roleInfo.Name, roleInfo.Description, accession.Normalize(roleInfo.Accession))

		c.Roles = append(c.Roles, r)
	}

	return c, nil
}

// 验证是否有这些权限
func (c *Controller) Require(a []accession.Accession) bool {
	for _, v := range a {
		if c.Has(v) {
			return true
		}
	}
	return false
}

// 检验是否拥有单独的权限
func (c *Controller) Has(a accession.Accession) bool {
	for _, r := range c.Roles {
		for _, v := range r.Accession {
			if v.Name == a.Name {
				return true
			}
		}
	}
	return false
}

// 根据 RBAC 鉴权的中间件
func Require(accesions ...accession.Accession) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			err error
			uid = ctx.GetString("uid") // 这个中间件必须安排在JWT的中间件后面, 所以这里是拿的到 UID 的
			c   *Controller
		)

		defer func() {
			if err != nil {
				ctx.JSON(http.StatusOK, schema.Response{
					Message: err.Error(),
					Data:    nil,
				})
				ctx.Abort()
			}
		}()

		if uid == "" {
			err = common_error.ErrNoPermission
		}

		if c, err = New(uid); err != nil {
			return
		}

		if c.Require(accesions) == false {
			err = common_error.ErrNoPermission
		}
	}
}
