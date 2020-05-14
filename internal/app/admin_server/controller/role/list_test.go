// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package role_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/role"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/rbac/accession"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	context := helper.Context{
		Uid: adminInfo.Id,
	}

	{
		var (
			name        = "vip"
			description = "VIP 用户"
			accessions  = accession.Stringify(accession.ProfileUpdate)
			n           = schema.Role{}
		)

		r := role.Create(helper.Context{
			Uid: adminInfo.Id,
		}, role.CreateParams{
			Name:        name,
			Description: description,
			Accession:   accessions,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer role.DeleteRoleByName(name)

		assert.Nil(t, r.Decode(&n))

		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}

	// 获取列表
	{
		r := role.GetList(context, role.Query{})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		list := make([]schema.Role, 0)

		assert.Nil(t, r.Decode(&list))

		assert.Equal(t, schema.DefaultLimit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.IsType(t, 1, r.Meta.Num)
		assert.IsType(t, int64(1), r.Meta.Total)

		if !assert.True(t, len(list) >= 1) {
			return
		}

		for _, n := range list {
			assert.IsType(t, "string", n.Name)
			assert.IsType(t, "string", n.Description)
			assert.IsType(t, []string{}, n.Accession)
			assert.IsType(t, true, n.BuildIn)
		}
	}
}

func TestGetListRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	{
		var (
			name        = "vip"
			description = "VIP 用户"
			accessions  = accession.Stringify(accession.ProfileUpdate)
			n           = schema.Role{}
		)

		r := role.Create(helper.Context{
			Uid: adminInfo.Id,
		}, role.CreateParams{
			Name:        name,
			Description: description,
			Accession:   accessions,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer role.DeleteRoleByName(name)

		assert.Nil(t, r.Decode(&n))

		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}

	{
		r := tester.HttpAdmin.Get("/v1/role", nil, &header)

		res := schema.Response{}

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		list := make([]schema.Role, 0)

		assert.Nil(t, res.Decode(&list))

		for _, n := range list {
			assert.IsType(t, "string", n.Name)
			assert.IsType(t, "string", n.Description)
			assert.IsType(t, []string{}, n.Accession)
			assert.IsType(t, true, n.BuildIn)
		}
	}
}
