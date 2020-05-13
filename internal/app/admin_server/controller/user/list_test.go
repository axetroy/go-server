// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package user_test

import (
	"encoding/json"
	user2 "github.com/axetroy/go-server/internal/app/admin_server/controller/user"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	context := helper.Context{
		Uid: adminInfo.Id,
	}

	// 获取列表
	{
		r := user2.GetList(context, user2.Query{})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		users := make([]schema.Profile, 0)

		assert.Nil(t, tester.Decode(r.Data, &users))

		assert.Equal(t, schema.DefaultLimit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)

		if !assert.True(t, len(users) >= 1) {
			return
		}

		for _, b := range users {
			assert.IsType(t, "string", b.Username)
			assert.IsType(t, "string", b.Id)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}

func TestGetListRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	{
		r := tester.HttpAdmin.Get("/v1/user", nil, &header)

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		users := make([]schema.Profile, 0)

		assert.Nil(t, tester.Decode(res.Data, &users))

		assert.Equal(t, schema.DefaultLimit, res.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, res.Meta.Page)

		if !assert.True(t, len(users) >= 1) {
			return
		}

		for _, b := range users {
			assert.IsType(t, "string", b.Username)
			assert.IsType(t, "string", b.Id)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}
