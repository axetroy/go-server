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

func TestCreate(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	// 创建
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		c := model.ConfigFieldSMTP{
			Server:    "0.0.0.0",
			Port:      10086,
			Username:  "username",
			Password:  "password",
			FromName:  "axetroy",
			FromEmail: "axetroy.dev@gmail.com",
		}

		body, _ := json.Marshal(c)

		r := tester.HttpAdmin.Post("/v1/config/smtp", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, schema.StatusSuccess, res.Status)

		defer database.DeleteRowByTable("config", "name", model.ConfigFieldNameSMTP.Field)

		n := schema.Config{}

		assert.Nil(t, res.Decode(&n))

		resultByte, err := json.Marshal(n.Fields)

		assert.Nil(t, err)

		resultConfig := model.ConfigFieldSMTP{}

		assert.Nil(t, json.Unmarshal(resultByte, &resultConfig))

		assert.Equal(t, c.Server, resultConfig.Server)
		assert.Equal(t, c.Port, resultConfig.Port)
		assert.Equal(t, c.Username, resultConfig.Username)
		assert.Equal(t, c.Password, resultConfig.Password)
		assert.Equal(t, c.FromName, resultConfig.FromName)
		assert.Equal(t, c.FromEmail, resultConfig.FromEmail)
	}
}
