// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package notification_test

import (
	"encoding/json"
	notificationAdmin "github.com/axetroy/go-server/internal/app/admin_server/controller/notification"
	"github.com/axetroy/go-server/internal/app/user_server/controller/notification"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarkRead(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	context := helper.Context{
		Uid: adminInfo.Id,
	}

	var testNotification schema.Notification

	// 创建一篇系统通知
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := notificationAdmin.Create(context, notificationAdmin.CreateParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testNotification = schema.Notification{}

		assert.Nil(t, r.Decode(&testNotification))

		defer notificationAdmin.DeleteNotificationById(testNotification.Id)

		assert.Equal(t, title, testNotification.Title)
		assert.Equal(t, content, testNotification.Content)
	}

	{
		// 不存在的用户标记系统通知为已读
		r := notification.MarkRead(helper.Context{
			Uid: "123123",
		}, testNotification.Id)

		assert.Equal(t, exception.UserNotExist.Code(), r.Status)
		assert.Equal(t, exception.UserNotExist.Error(), r.Message)
	}

	{
		// 用测试用户标记为已读
		r := notification.MarkRead(helper.Context{
			Uid: userInfo.Id,
		}, testNotification.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.Equal(t, nil, r.Data)

		defer notificationAdmin.DeleteNotificationMarkById(testNotification.Id)
	}

	{
		// 再读取这条系统通知
		notificationMarkInfo := model.NotificationMark{}

		assert.Nil(t, database.Db.Where("id = ?", testNotification.Id).Last(&notificationMarkInfo).Error)
		assert.Equal(t, userInfo.Id, notificationMarkInfo.Uid)
	}

	{
		// 获取详情
		r := notification.Get(context, testNotification.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		n := schema.Notification{}
		assert.Nil(t, r.Decode(&n))

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

	defer tester.DeleteUserByUserName(userInfo.Username)

	context := helper.Context{
		Uid: adminInfo.Id,
	}

	var testNotification schema.Notification

	// 创建一篇系统通知
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := notificationAdmin.Create(context, notificationAdmin.CreateParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testNotification = schema.Notification{}

		assert.Nil(t, r.Decode(&testNotification))

		notificationId = testNotification.Id

		defer notificationAdmin.DeleteNotificationById(testNotification.Id)

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

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		assert.Equal(t, nil, res.Data)

		// 再读取这条系统通知
		notificationMarkInfo := model.NotificationMark{}

		assert.Nil(t, database.Db.Where("id = ?", testNotification.Id).Last(&notificationMarkInfo).Error)
		assert.Equal(t, userInfo.Id, notificationMarkInfo.Uid)
	}
}
