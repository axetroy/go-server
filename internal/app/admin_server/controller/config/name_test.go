// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetNameRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	// 获取列表
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		r := tester.HttpAdmin.Get("/v1/config/name", nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, schema.StatusSuccess, res.Status)

		var list []schema.Name

		assert.Nil(t, res.Decode(&list))

		resultByte, err := json.Marshal(list)

		assert.Nil(t, err)

		expectByte, err := json.Marshal(model.ConfigFields)

		assert.Nil(t, err)

		assert.Equal(t, expectByte, resultByte)
	}
}
