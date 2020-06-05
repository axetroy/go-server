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
	"testing"
)

func TestGetList(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	adminRes := admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "test",
		Password: "test123",
		Name:     "test",
	}, false)

	newAdmin := schema.AdminProfile{}

	assert.Nil(t, adminRes.Decode(&newAdmin))

	defer admin.DeleteAdminByAccount("test")

	r := admin.GetList(helper.Context{Uid: adminInfo.Id}, admin.Query{})

	list := make([]schema.AdminProfile, 0)

	assert.Nil(t, r.Decode(&list))

	assert.Equal(t, schema.DefaultLimit, r.Meta.Limit)
	assert.Equal(t, schema.DefaultPage, r.Meta.Page)
	assert.IsType(t, 1, r.Meta.Num)
	assert.IsType(t, int64(1), r.Meta.Total)

	assert.True(t, len(list) >= 1)

	for _, b := range list {
		assert.IsType(t, "string", b.Id)
		assert.IsType(t, "string", b.Username)
		assert.IsType(t, "string", b.CreatedAt)
		assert.IsType(t, "string", b.UpdatedAt)
	}
}

func TestGetListRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	adminRes := admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "test",
		Password: "test123",
		Name:     "test",
	}, false)

	newAdmin := schema.AdminProfile{}

	assert.Nil(t, adminRes.Decode(&newAdmin))

	defer admin.DeleteAdminByAccount("test")

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	r := tester.HttpAdmin.Get("/v1/admin", nil, &header)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	list := make([]schema.AdminProfile, 0)

	assert.Nil(t, res.Decode(&list))

	for _, b := range list {
		assert.IsType(t, "string", b.Id)
		assert.IsType(t, "string", b.Username)
		assert.IsType(t, "string", b.Name)
		assert.IsType(t, "string", b.CreatedAt)
		assert.IsType(t, "string", b.UpdatedAt)
	}
}
