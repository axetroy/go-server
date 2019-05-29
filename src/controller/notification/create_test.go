package notification_test

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/notification"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	var (
		adminUid string
	)
	// 先登陆获取管理员的Token
	{
		// 登陆超级管理员-成功

		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "admin",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := schema.AdminProfileWithToken{}

		if err := tester.Decode(r.Data, &adminInfo); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, "admin", adminInfo.Username)
		assert.True(t, len(adminInfo.Token) > 0)

		if c, er := token.Parse(token.Prefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			adminUid = c.Uid
		}
	}

	// 创建一篇系统通知
	{
		var (
			title   = "TestCreate"
			content = "TestCreate"
		)

		r := notification.Create(controller.Context{
			Uid: adminUid,
		}, notification.CreateParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Notification{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer notification.DeleteNotificationById(n.Id)

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
	}
}
