// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	adminRes := admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "test",
		Password: "test",
		Name:     "test",
	}, false)

	newAdmin := admin_schema.AdminProfile{}

	assert.Nil(t, tester.Decode(adminRes.Data, &newAdmin))

	defer admin.DeleteAdminByAccount("test")

	r := admin.GetList(schema.Context{Uid: adminInfo.Id}, admin.Query{})

	list := make([]admin_schema.AdminProfile, 0)

	assert.Nil(t, tester.Decode(r.Data, &list))

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
		Password: "test",
		Name:     "test",
	}, false)

	newAdmin := admin_schema.AdminProfile{}

	assert.Nil(t, tester.Decode(adminRes.Data, &newAdmin))

	defer admin.DeleteAdminByAccount("test")

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	r := tester.HttpAdmin.Get("/v1/admin", nil, &header)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	list := make([]admin_schema.AdminProfile, 0)

	assert.Nil(t, tester.Decode(res.Data, &list))

	for _, b := range list {
		assert.IsType(t, "string", b.Id)
		assert.IsType(t, "string", b.Username)
		assert.IsType(t, "string", b.Name)
		assert.IsType(t, "string", b.CreatedAt)
		assert.IsType(t, "string", b.UpdatedAt)
	}
}
