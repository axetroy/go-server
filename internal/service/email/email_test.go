// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package email_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/email"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetMailerConfig(t *testing.T) {
	{
		c, err := email.GetMailerConfig()

		assert.Nil(t, c)
		assert.NotNil(t, err)
	}

	{
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

		result, err := email.GetMailerConfig()

		assert.NotNil(t, c)
		assert.Nil(t, err)

		assert.Equal(t, c.Server, result.Server)
		assert.Equal(t, c.Port, result.Port)
		assert.Equal(t, c.Username, result.Username)
		assert.Equal(t, c.Password, result.Password)
		assert.Equal(t, c.FromName, result.FromName)
		assert.Equal(t, c.FromEmail, result.FromEmail)
	}
}
