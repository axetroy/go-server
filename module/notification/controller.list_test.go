// Copyright 2019 Axetroy. All rights reserved. MIT license.
package notification_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_schema"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/notification"
	"github.com/axetroy/go-server/module/notification/notification_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	{
		var (
			adminUid string
		)
		// 1. 先登陆获取管理员的Token
		{
			r := admin.Login(admin.SignInParams{
				Username: "admin",
				Password: "admin",
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			adminInfo := admin_schema.AdminProfileWithToken{}

			assert.Nil(t, tester.Decode(r.Data, &adminInfo))

			if c, er := token.Parse(token.Prefix+" "+adminInfo.Token, true); er != nil {
				t.Error(er)
			} else {
				adminUid = c.Uid
			}
		}

		// 2. 先创建一个通知作为测试
		{
			var (
				title   = "TestGetList"
				content = "TestGetList"
			)

			r := notification.Create(schema.Context{
				Uid: adminUid,
			}, notification.CreateParams{
				Title:   title,
				Content: content,
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			n := notification_schema.Notification{}

			assert.Nil(t, tester.Decode(r.Data, &n))

			defer notification.DeleteNotificationById(n.Id)
		}

		// 3. 获取列表
		{
			var (
				data = make([]notification_schema.Notification, 0)
			)
			query := schema.Query{
				Limit: 20,
			}
			r := notification.GetListByUser(schema.Context{}, notification.Query{
				Query: query,
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			assert.Nil(t, tester.Decode(r.Data, &data))
			assert.Equal(t, query.Limit, r.Meta.Limit)
			assert.Equal(t, schema.DefaultPage, r.Meta.Page)
			assert.IsType(t, 1, r.Meta.Num)
			assert.IsType(t, int64(1), r.Meta.Total)

			assert.True(t, len(data) > 0)
		}
	}
}

func TestGetListRouter(t *testing.T) {
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
			title   = "test"
			content = "test"
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

	// 管理员接口获取
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		r := tester.HttpAdmin.Get("/v1/notification", nil, &header)
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

		banners := make([]notification_schema.Notification, 0)

		assert.Nil(t, tester.Decode(res.Data, &banners))

		for _, b := range banners {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}

	// 普通用户获取通知
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + userInfo.Token,
		}

		r := tester.HttpUser.Get("/v1/notification", nil, &header)
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

		list := make([]notification_schema.Notification, 0)

		assert.Nil(t, tester.Decode(res.Data, &list))

		for _, b := range list {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}
