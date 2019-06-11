// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin_test

import (
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteAdminByAccount(t *testing.T) {
	{
		// 创建已存在的管理员
		r := admin.CreateAdmin(admin.CreateAdminParams{
			Account:  "admin123",
			Name:     "test",
			Password: "123",
		}, false)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
	}

	{
		// 删除管理员
		admin.DeleteAdminByAccount("admin123")
	}

	{
		// 获取管理员信息
		adminInfo := admin_model.Admin{
			Username: "admin123",
		}

		err := database.Db.Where(&adminInfo).First(&adminInfo).Error

		assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
	}
}
