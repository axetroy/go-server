// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package system_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/system"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetSystemInfo(t *testing.T) {
	r := system.GetSystemInfo()

	assert.Equal(t, 1, r.Status)
	assert.Equal(t, "", r.Message)
}

func TestGetSystemInfoRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	r := tester.HttpAdmin.Get("/v1/system", nil, &header)

	res := schema.Response{}

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

	n := system.Info{}

	assert.Nil(t, tester.Decode(res.Data, &n))
}
