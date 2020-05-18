// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package notification_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/admin"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/notification"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNotificationListByAdmin(t *testing.T) {
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

			adminInfo := schema.AdminProfileWithToken{}

			assert.Nil(t, r.Decode(&adminInfo))

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

			r := notification.Create(helper.Context{
				Uid: adminUid,
			}, notification.CreateParams{
				Title:   title,
				Content: content,
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			n := schema.Notification{}

			assert.Nil(t, r.Decode(&n))

			defer notification.DeleteNotificationById(n.Id)
		}

		// 3. 获取列表
		{
			var (
				data = make([]schema.Notification, 0)
			)
			query := schema.Query{
				Limit: 20,
			}
			r := notification.GetNotificationListByAdmin(helper.Context{}, notification.Query{
				Query: query,
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			assert.Nil(t, r.Decode(&data))
			assert.Equal(t, query.Limit, r.Meta.Limit)
			assert.Equal(t, schema.DefaultPage, r.Meta.Page)
			assert.IsType(t, 1, r.Meta.Num)
			assert.IsType(t, int64(1), r.Meta.Total)

			assert.True(t, len(data) > 0)
		}
	}
}

func TestGetNotificationListByAdminRouter(t *testing.T) {
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

		assert.Nil(t, r.Decode(&testNotification))

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

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, "", res.Message)
		assert.Equal(t, schema.StatusSuccess, res.Status)

		banners := make([]schema.Notification, 0)

		assert.Nil(t, res.Decode(&banners))

		for _, b := range banners {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}
