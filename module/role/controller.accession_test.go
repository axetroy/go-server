// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/role"
	"github.com/axetroy/go-server/rbac/accession"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetAccession(t *testing.T) {
	{
		r := role.Get("123123")

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, role.ErrRoleNotExist.Error(), r.Message)
	}

	r := role.GetAccession()

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	accessions := make([]*accession.Accession, 0)

	assert.Nil(t, tester.Decode(r.Data, &accessions))
	assert.Equal(t, accessions, accession.List)
}

func TestGetAccessionRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	r := tester.HttpAdmin.Get("/v1/role/accession", nil, &header)
	res := schema.Response{}

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	accessions := make([]*accession.Accession, 0)

	assert.Nil(t, tester.Decode(res.Data, &accessions))
	assert.Equal(t, accessions, accession.List)
}
