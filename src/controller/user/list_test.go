// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/user"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	// 获取列表
	{
		r := user.GetList(context, user.Query{})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		users := make([]schema.Profile, 0)

		assert.Nil(t, tester.Decode(r.Data, &users))

		assert.Equal(t, schema.DefaultLimit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.IsType(t, 1, r.Meta.Num)
		assert.IsType(t, int64(1), r.Meta.Total)

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

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	{
		r := tester.HttpAdmin.Get("/v1/user", nil, &header)

		res := schema.List{}

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
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
		assert.IsType(t, 1, res.Meta.Num)
		assert.IsType(t, int64(1), res.Meta.Total)

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
