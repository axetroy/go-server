// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package notification_test

import (
	"github.com/axetroy/go-server/internal/app/admin_server/controller/notification"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	adminInfo, err := tester.LoginAdmin()

	assert.Nil(t, err)

	// 创建一篇系统通知
	{
		var (
			title   = "TestCreate"
			content = "TestCreate"
		)

		r := notification.Create(helper.Context{
			Uid: adminInfo.Id,
		}, notification.CreateParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Notification{}

		assert.Nil(t, r.Decode(&n))

		defer notification.DeleteNotificationById(n.Id)

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
	}
}
