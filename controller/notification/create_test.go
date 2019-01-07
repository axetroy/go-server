package notification_test

import (
	"fmt"
	"github.com/axetroy/go-server/controller"
	"github.com/axetroy/go-server/controller/admin"
	"github.com/axetroy/go-server/controller/notification"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/go-server/token"
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

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := admin.SignInResponse{}

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
			title   = "test"
			content = "test"
		)

		r := notification.Create(controller.Context{
			Uid: adminUid,
		}, notification.CreateParams{
			Tittle:  title,
			Content: content,
		})

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := notification.Notification{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer func() {
			notification.DeleteNotificationById(n.Id)
		}()

		assert.Equal(t, title, n.Tittle)
		assert.Equal(t, content, n.Content)

		fmt.Printf("%+v\n", n)
	}
}
