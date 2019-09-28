// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/role"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/rbac/accession"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	{
		r := role.Get("123123")

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.RoleNotExist.Error(), r.Message)
	}

	adminInfo, _ := tester.LoginAdmin()

	var (
		name        = "vip"
		description = "VIP 用户"
		accessions  = accession.Stringify(accession.ProfileUpdate)
		n           = schema.Role{}
	)

	{

		r := role.Create(controller.Context{
			Uid: adminInfo.Id,
		}, role.CreateParams{
			Name:        name,
			Description: description,
			Accession:   accessions,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer role.DeleteRoleByName(name)

		assert.Nil(t, tester.Decode(r.Data, &n))

		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}

	{
		r := role.Get(n.Name)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		roleInfo := schema.Role{}

		assert.Nil(t, tester.Decode(r.Data, &roleInfo))
		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}
}

func TestGetRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	var (
		name        = "vip"
		description = "VIP 用户"
		accessions  = accession.Stringify(accession.ProfileUpdate)
		n           = schema.Role{}
	)

	{

		r := role.Create(controller.Context{
			Uid: adminInfo.Id,
		}, role.CreateParams{
			Name:        name,
			Description: description,
			Accession:   accessions,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer role.DeleteRoleByName(name)

		assert.Nil(t, tester.Decode(r.Data, &n))

		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}

	// 获取详情
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		r := tester.HttpAdmin.Get("/v1/role/r/"+n.Name, nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		d := schema.Role{}

		assert.Nil(t, tester.Decode(res.Data, &d))

		assert.Equal(t, n.Name, d.Name)
		assert.Equal(t, n.Description, d.Description)
		assert.Equal(t, n.Accession, d.Accession)
		assert.Equal(t, n.Note, d.Note)
	}
}
