// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/core/controller"
	"github.com/axetroy/go-server/core/controller/role"
	"github.com/axetroy/go-server/core/rbac/accession"
	"github.com/axetroy/go-server/core/schema"
	"github.com/axetroy/go-server/core/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	var (
		name        = "vip"
		description = "VIP 用户"
		accessions  = accession.Stringify(accession.ProfileUpdate)
	)

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

	n := schema.Role{}

	assert.Nil(t, tester.Decode(r.Data, &n))

	assert.Equal(t, name, n.Name)
	assert.Equal(t, description, n.Description)
	assert.Equal(t, accessions, n.Accession)
	assert.Equal(t, false, n.BuildIn)
}

func TestCreateRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	// 创建 banner
	{
		var (
			name        = "vip"
			description = "VIP 用户"
			accessions  = accession.Stringify(accession.ProfileUpdate)
		)

		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		body, _ := json.Marshal(&role.CreateParams{
			Name:        name,
			Description: description,
			Accession:   accessions,
		})

		r := tester.HttpAdmin.Post("/v1/role", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		n := schema.Role{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		defer role.DeleteRoleByName(n.Name)

		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}
}
