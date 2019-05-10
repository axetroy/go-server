package src_test

import (
	"github.com/axetroy/go-server/src"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	src.Init()

	// 删除管理员
	defer admin.DeleteAdminByAccount("admin")

	// 获取管理员
	adminInfo := model.Admin{
		Username: "admin",
	}

	assert.Nil(t, service.Db.Where(&adminInfo).First(&adminInfo).Error)
}
