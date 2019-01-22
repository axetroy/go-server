package admin_test

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAdminInfo(t *testing.T) {
	var adminUid string
	{
		// 1. 创建一个测试管理员
		input := admin.CreateAdminParams{
			Account:  "test",
			Name:     "test",
			Password: "123",
		}

		r := admin.CreateAdmin(input, false)

		assert.Equal(t, r.Status, schema.StatusSuccess)
		assert.Equal(t, r.Message, "")

		defer func() {
			// 删除这个刚创建的管理员
			admin.DeleteAdminByAccount(input.Account)
		}()

		detail := schema.AdminProfile{}

		if err := tester.Decode(r.Data, &detail); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, input.Account, detail.Username)
		assert.Equal(t, input.Name, detail.Name)

		adminUid = detail.Id

	}

	{
		// 2. 获取管理员信息

		r := admin.GetAdminInfo(controller.Context{
			Uid: adminUid,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		detail := schema.AdminProfile{}

		assert.Nil(t, tester.Decode(r.Data, &detail))

		assert.Equal(t, adminUid, detail.Id)
	}
}
