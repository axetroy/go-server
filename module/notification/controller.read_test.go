// Copyright 2019 Axetroy. All rights reserved. MIT license.
package notification_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/notification"
	"github.com/axetroy/go-server/module/notification/notification_model"
	"github.com/axetroy/go-server/module/notification/notification_schema"
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarkRead(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	ctx := schema.Context{
		Uid: adminInfo.Id,
	}

	var testNotification notification_schema.Notification

	// 创建一篇系统通知
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := notification.Create(ctx, notification.CreateParams{
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

	{
		// 不存在的用户标记系统通知为已读
		r := notification.MarkRead(schema.Context{
			Uid: "123123",
		}, testNotification.Id)

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, user_error.ErrUserNotExist.Error(), r.Message)
	}

	{
		// 用测试用户标记为已读
		r := notification.MarkRead(schema.Context{
			Uid: userInfo.Id,
		}, testNotification.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.True(t, r.Data.(bool))

		defer notification.DeleteNotificationMarkById(testNotification.Id)
	}

	{
		// 再读取这条系统通知
		notificationMarkInfo := notification_model.NotificationMark{}

		assert.Nil(t, database.Db.Where("id = ?", testNotification.Id).Last(&notificationMarkInfo).Error)
		assert.Equal(t, userInfo.Id, notificationMarkInfo.Uid)
	}

	{
		// 获取详情
		r := notification.Get(ctx, testNotification.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		n := notification_schema.Notification{}
		assert.Nil(t, tester.Decode(r.Data, &n))

		assert.Equal(t, n.Id, testNotification.Id)
		assert.Equal(t, n.Title, testNotification.Title)
		assert.Equal(t, n.Content, testNotification.Content)
		assert.Equal(t, false, testNotification.Read)
	}
}

func TestReadRouter(t *testing.T) {
	var notificationId string
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := schema.Context{
		Uid: adminInfo.Id,
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

		notificationId = testNotification.Id

		defer notification.DeleteNotificationById(testNotification.Id)

		assert.Equal(t, title, testNotification.Title)
		assert.Equal(t, content, testNotification.Content)
	}

	// 标记为已读
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + userInfo.Token,
		}

		r := tester.HttpUser.Put("/v1/notification/n/"+notificationId+"/read", nil, &header)
		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		assert.True(t, res.Data.(bool))

		// 再读取这条系统通知
		notificationMarkInfo := notification_model.NotificationMark{}

		assert.Nil(t, database.Db.Where("id = ?", testNotification.Id).Last(&notificationMarkInfo).Error)
		assert.Equal(t, userInfo.Id, notificationMarkInfo.Uid)
	}
}
