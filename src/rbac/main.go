package rbac

import (
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/rbac/accession"
	"github.com/axetroy/go-server/src/rbac/role"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/jinzhu/gorm"
)

type Controller struct {
	Roles []*role.Role
}

func New(uid string) (c *Controller, err error) {
	c = &Controller{}

	userInfo := model.User{
		Id: uid,
	}

	if err = database.Db.First(&userInfo).Error; err != nil {
		return
	}

	if len(userInfo.Role) == 0 {
		err = exception.NoPermission
		return
	}

	for _, roleName := range userInfo.Role {
		roleInfo := model.Role{
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
