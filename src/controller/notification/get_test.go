// Copyright 2019 Axetroy. All rights reserved. MIT license.
package notification_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/notification"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	var testNotification schema.Notification

	// 创建一篇系统通知
	{
		var (
			title   = "TestGet"
			content = "TestGet"
		)

		r := notification.Create(context, notification.CreateParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testNotification = schema.Notification{}

		assert.Nil(t, tester.Decode(r.Data, &testNotification))

		defer notification.DeleteNotificationById(testNotification.Id)

		assert.Equal(t, title, testNotification.Title)
		assert.Equal(t, content, testNotification.Content)
	}

	// 获取详情
	{
		r := notification.Get(context, testNotification.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		n := schema.Notification{}
		assert.Nil(t, tester.Decode(r.Data, &n))

		assert.Equal(t, n.Id, testNotification.Id)
		assert.Equal(t, n.Title, testNotification.Title)
		assert.Equal(t, n.Content, testNotification.Content)
		assert.Equal(t, false, testNotification.Read)
	}
}

func TestGetRouter(t *testing.T) {
	var notificationId string
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	var testNotification schema.Notification

	// 创建一篇系统通知
	{
		var (
			title   = "test"
			content = "test"
		)

		r := notification.Create(context, notification.CreateParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testNotification = schema.Notification{}

		assert.Nil(t, tester.Decode(r.Data, &testNotification))

		notificationId = testNotification.Id

		defer notification.DeleteNotificationById(testNotification.Id)

		assert.Equal(t, title, testNotification.Title)
		assert.Equal(t, content, testNotification.Content)
	}

	// 管理员接口获取
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		r := tester.HttpAdmin.Get("/v1/notification/n/"+notificationId, nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		n := schema.Notification{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		assert.Equal(t, "test", n.Title)
		assert.Equal(t, "test", n.Content)
		assert.IsType(t, "string", n.CreatedAt)
		assert.IsType(t, "string", n.UpdatedAt)
	}

	// 普通用户获取通知
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + userInfo.Token,
		}

		r := tester.HttpUser.Get("/v1/notification/n/"+notificationId, nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		n := schema.Notification{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		assert.Equal(t, "test", n.Title)
		assert.Equal(t, "test", n.Content)
		assert.IsType(t, "string", n.CreatedAt)
		assert.IsType(t, "string", n.UpdatedAt)
	}
}
