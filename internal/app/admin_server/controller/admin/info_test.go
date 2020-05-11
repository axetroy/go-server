// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/admin"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetAdminInfo(t *testing.T) {
	var adminUid string
	{
		// 1. 创建一个测试管理员
		input := admin.CreateAdminParams{
			Account:  "test",
			Name:     "test",
			Password: "123",
		}

		r := admin.CreateAdmin(input, false)

		assert.Equal(t, r.Status, schema.StatusSuccess)
		assert.Equal(t, r.Message, "")

		defer admin.DeleteAdminByAccount(input.Account)

		detail := schema.AdminProfile{}

		if err := tester.Decode(r.Data, &detail); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, input.Account, detail.Username)
		assert.Equal(t, input.Name, detail.Name)

		adminUid = detail.Id

	}

	{
		// 2. 获取管理员信息

		r := admin.GetAdminInfo(helper.Context{
			Uid: adminUid,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		detail := schema.AdminProfile{}

		assert.Nil(t, tester.Decode(r.Data, &detail))

		assert.Equal(t, adminUid, detail.Id)
	}
}

func TestGetAdminInfoRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	r := tester.HttpAdmin.Get("/v1/profile", nil, &header)
	res := schema.Response{}

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

	n := schema.AdminProfile{}

	assert.Nil(t, tester.Decode(res.Data, &n))

	assert.Equal(t, adminInfo.Name, n.Name)
	assert.Equal(t, adminInfo.Username, n.Username)
	assert.Equal(t, adminInfo.Status, n.Status)
}

func TestGetAdminInfoById(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	var adminUid string
	{
		// 1. 创建一个测试管理员
		input := admin.CreateAdminParams{
			Account:  "123123",
			Name:     "123123",
			Password: "123123",
		}

		r := admin.CreateAdmin(input, false)

		assert.Equal(t, r.Status, schema.StatusSuccess)
		assert.Equal(t, r.Message, "")

		defer admin.DeleteAdminByAccount(input.Account)

		detail := schema.AdminProfile{}

		if err := tester.Decode(r.Data, &detail); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, input.Account, detail.Username)
		assert.Equal(t, input.Name, detail.Name)

		adminUid = detail.Id

	}

	{
		// 2. 获取管理员信息

		r := admin.GetAdminInfoById(helper.Context{
			Uid: adminInfo.Id,
		}, adminUid)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		detail := schema.AdminProfile{}

		assert.Nil(t, tester.Decode(r.Data, &detail))

		assert.Equal(t, adminUid, detail.Id)
		assert.Equal(t, "123123", detail.Username)
		assert.Equal(t, "123123", detail.Name)
	}
}

func TestGetAdminInfoByIdRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	r := tester.HttpAdmin.Get("/v1/admin/a/"+adminInfo.Id, nil, &header)
	res := schema.Response{}

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

	n := schema.AdminProfile{}

	assert.Nil(t, tester.Decode(res.Data, &n))

	assert.Equal(t, adminInfo.Name, n.Name)
	assert.Equal(t, adminInfo.Username, n.Username)
	assert.Equal(t, adminInfo.Status, n.Status)
}
