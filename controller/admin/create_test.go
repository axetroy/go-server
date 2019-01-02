package admin_test

import (
	"github.com/axetroy/go-server/controller/admin"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAdmin(t *testing.T) {
	// 不能创建超级管理员，因为只能存在一个超级管理员
	{
		r := admin.CreateAdmin(admin.CreateAdminParams{
			Account:  "123123",
			Name:     "test",
			Password: "123",
		}, true)

		assert.Equal(t, r.Status, response.StatusFail)
	}

	// 创建已存在的管理员
	{
		r := admin.CreateAdmin(admin.CreateAdminParams{
			Account:  "admin",
			Name:     "test",
			Password: "123",
		}, true)

		assert.Equal(t, r.Status, response.StatusFail)
		assert.Equal(t, r.Message, exception.AdminExist.Error())
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
	// TODO: 测试路由是否正常工作
}
