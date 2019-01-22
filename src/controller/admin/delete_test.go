package admin_test

import (
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteAdminByAccount(t *testing.T) {
	{
		// 创建已存在的管理员
		r := admin.CreateAdmin(admin.CreateAdminParams{
			Account:  "admin123",
			Name:     "test",
			Password: "123",
		}, false)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
	}

	{
		// 删除管理员
		admin.DeleteAdminByAccount("admin123")
	}

	{
		// TODO: 获取管理员信息
	}
}
