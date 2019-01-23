package notification_test

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/notification"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestMarkRead(t *testing.T) {
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

		if c, er := util.ParseToken(util.TokenPrefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			adminUid = c.Uid
		}
	}

	context := controller.Context{
		Uid: adminUid,
	}

	var testNotification schema.Notification

	// 创建一篇系统通知
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := notification.Create(context, notification.CreateParams{
			Tittle:  title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testNotification = schema.Notification{}

		assert.Nil(t, tester.Decode(r.Data, &testNotification))

		defer notification.DeleteNotificationById(testNotification.Id)

		assert.Equal(t, title, testNotification.Tittle)
		assert.Equal(t, content, testNotification.Content)
	}

	{
		// 不存在的用户标记系统通知为已读
		r := notification.MarkRead(controller.Context{
			Uid: "123123",
		}, testNotification.Id)

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.UserNotExist.Error(), r.Message)
	}

	var testUser schema.Profile

	{
		// 创建一个测试用户
		// 1。 创建测试账号
		rand.Seed(111)
		username := "test-TestMarkRead"
		password := "123123"

		r := auth.SignUp(auth.SignUpParams{
			Username: &username,
			Password: password,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = schema.Profile{}

		if err := tester.Decode(r.Data, &testUser); err != nil {
			t.Error(err)
			return
		}

		defer auth.DeleteUserByUserName(username)
	}

	{
		// 用测试用户标记为已读
		r := notification.MarkRead(controller.Context{
			Uid: testUser.Id,
		}, testNotification.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.True(t, r.Data.(bool))

		defer notification.DeleteNotificationMarkById(testNotification.Id)
	}

	{
		// 再读取这条系统通知
		notificationMarkInfo := model.NotificationMark{}

		assert.Nil(t, service.Db.Where("id = ?", testNotification.Id).Last(&notificationMarkInfo).Error)
		assert.Equal(t, testUser.Id, notificationMarkInfo.Uid)
	}

	{
		// TODO: 用户再获取这个通知，会标记为已读
	}
}

func TestReadRouter(t *testing.T) {
	// TODO: 完善HTTP的测试用例
}
