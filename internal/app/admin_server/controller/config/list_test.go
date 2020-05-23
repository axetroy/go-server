// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetListRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	c := model.ConfigFieldSMTP{
		Server:    "0.0.0.0",
		Port:      10086,
		Username:  "username",
		Password:  "password",
		FromName:  "axetroy",
		FromEmail: "axetroy.dev@gmail.com",
	}

	// 创建
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		body, _ := json.Marshal(c)

		r := tester.HttpAdmin.Post("/v1/config/smtp", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, schema.StatusSuccess, res.Status)

		defer database.DeleteRowByTable("config", "name", model.ConfigFieldNameSMTP.Field)
	}

	// 获取列表
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		r := tester.HttpAdmin.Get("/v1/config", nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, schema.StatusSuccess, res.Status)

		n := make([]schema.Config, 0)

		assert.Nil(t, res.Decode(&n))

		assert.Equal(t, schema.DefaultLimit, res.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, res.Meta.Page)
		assert.IsType(t, 1, res.Meta.Num)
		assert.IsType(t, int64(1), res.Meta.Total)

		assert.True(t, len(n) >= 1)

		for _, b := range n {
			_, ok := b.Fields.(map[string]interface{})
			assert.True(t, ok)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}
