package admin_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/controller/admin"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/go-server/token"
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
	// 不能创建超级管理员，因为只能存在一个超级管理员
	{
		r := admin.CreateAdmin(admin.CreateAdminParams{
			Account:  "123123",
			Name:     "test",
			Password: "123",
		}, true)

		assert.Equal(t, response.StatusFail, r.Status)
	}

	// 创建已存在的管理员
	{
		r := admin.CreateAdmin(admin.CreateAdminParams{
			Account:  "admin",
			Name:     "test",
			Password: "123",
		}, true)

		assert.Equal(t, response.StatusFail, r.Status)
		assert.Equal(t, exception.AdminExist.Error(), r.Message)
	}

	// 创建普通的管理员成功
	{
		input := admin.CreateAdminParams{
			Account:  "test",
			Name:     "test",
			Password: "123",
		}

		r := admin.CreateAdmin(input, false)

		assert.Equal(t, r.Status, response.StatusSuccess)
		assert.Equal(t, r.Message, "")

		defer func() {
			// 删除这个刚创建的管理员
			admin.DeleteAdminByAccount(input.Account)
		}()

		detail := admin.Detail{}

		if err := tester.Decode(r.Data, &detail); err != nil {
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
		password := "12312"

		body, _ := json.Marshal(&admin.CreateAdminParams{
			Account:  username,
			Password: password,
			Name:     username,
		})

		r := tester.Http.Post("/v1/admin/admin/create", body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := response.Response{}

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
			return
		}

		if !assert.Equal(t, response.StatusFail, res.Status) {
			return
		}
		if !assert.Equal(t, exception.InvalidToken.Error(), res.Message) {
			return
		}
	}

	{
		// 拿正确的Token创建管理员
	}
}
