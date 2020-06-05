// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/admin"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func init() {
	admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "admin",
		Password: "admin",
		Name:     "admin",
	}, true)
}

func TestCreateAdmin(t *testing.T) {
	// 创建已存在的管理员
	{
		r := admin.CreateAdmin(admin.CreateAdminParams{
			Account:  "admin",
			Name:     "test",
			Password: "123123",
		}, true)

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.AdminExist.Error(), r.Message)
	}

	// 创建普通的管理员成功
	{
		input := admin.CreateAdminParams{
			Account:  "test",
			Name:     "test",
			Password: "123123",
		}

		r := admin.CreateAdmin(input, false)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer func() {
			// 删除这个刚创建的管理员
			admin.DeleteAdminByAccount(input.Account)
		}()

		detail := schema.AdminProfile{}

		if err := r.Decode(&detail); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, detail.Username, input.Account)
		assert.Equal(t, detail.Name, input.Name)
	}
}

func TestCreateAdminRouter(t *testing.T) {
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " 12312",
		}

		username := "test-TestCreateAdminRouter"
		password := "123123"

		body, _ := json.Marshal(&admin.CreateAdminParams{
			Account:  username,
			Password: password,
			Name:     username,
		})

		r := tester.HttpAdmin.Post("/v1/admin", body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, exception.InvalidToken.Code(), res.Status)
		assert.Equal(t, exception.InvalidToken.Error(), res.Message)
	}

	{
		// 拿正确的Token创建管理员
	}
}
