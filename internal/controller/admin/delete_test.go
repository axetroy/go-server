// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/controller"
	"github.com/axetroy/go-server/internal/controller/admin"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDeleteAdminByAccount(t *testing.T) {
	{
		// 创建管理员
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
		adminInfo := model.Admin{
			Username: "admin123",
		}

		err := database.Db.Where(&adminInfo).First(&adminInfo).Error

		assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
	}
}

func TestDeleteAdminById(t *testing.T) {
	adminInfo, err := tester.LoginAdmin()

	assert.Nil(t, err)

	var uid string

	{
		// 创建管理员
		r := admin.CreateAdmin(admin.CreateAdminParams{
			Account:  "admin321",
			Name:     "test",
			Password: "123",
		}, false)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer admin.DeleteAdminByAccount("admin321")

		detail := schema.AdminProfileWithToken{}

		assert.Nil(t, tester.Decode(r.Data, &detail))

		uid = detail.Id
	}

	{
		// 删除管理员
		r := admin.DeleteAdminById(controller.Context{Uid: adminInfo.Id}, uid)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
	}

	{
		// 获取管理员信息
		d := model.Admin{
			Id: uid,
		}

		err := database.Db.Where(&d).First(&d).Error

		assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
	}
}

func TestDeleteAdminByIdRouter(t *testing.T) {
	adminInfo, err := tester.LoginAdmin()

	assert.Nil(t, err)

	var uid string

	{
		// 创建管理员
		r := admin.CreateAdmin(admin.CreateAdminParams{
			Account:  "admin321",
			Name:     "test",
			Password: "123",
		}, false)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer admin.DeleteAdminByAccount("admin321")

		detail := schema.AdminProfileWithToken{}

		assert.Nil(t, tester.Decode(r.Data, &detail))

		uid = detail.Id
	}

	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		r := tester.HttpAdmin.Delete("/v1/admin/a/"+uid, nil, &header)

		assert.Equal(t, http.StatusOK, r.Code)

		res := schema.Response{}

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)
	}
}
