// Copyright 2019 Axetroy. All rights reserved. MIT license.
package notification_test

import (
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_schema"
	"github.com/axetroy/go-server/module/notification"
	"github.com/axetroy/go-server/module/notification/notification_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdate(t *testing.T) {
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

		adminInfo := admin_schema.AdminProfileWithToken{}

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

	context := schema.Context{
		Uid: adminUid,
	}

	var testNotification notification_schema.Notification

	// 创建一篇系统通知
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := notification.Create(context, notification.CreateParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testNotification = notification_schema.Notification{}

		assert.Nil(t, tester.Decode(r.Data, &testNotification))

		defer notification.DeleteNotificationById(testNotification.Id)

		assert.Equal(t, title, testNotification.Title)
		assert.Equal(t, content, testNotification.Content)
	}

	// 更新系统通知
	{
		var (
			newTittle  = "123123"
			newContent = "123123"
			newNote    = "123123"
		)

		r := notification.Update(context, testNotification.Id, notification.UpdateParams{
			Title:   &newTittle,
			Content: &newContent,
			Note:    &newNote,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := notification_schema.Notification{}
		assert.Nil(t, tester.Decode(r.Data, &n))
		assert.Equal(t, newTittle, n.Title)
		assert.Equal(t, newContent, n.Content)
		assert.Equal(t, newNote, *n.Note)
	}
}
