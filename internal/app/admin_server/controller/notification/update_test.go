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

func TestUpdate(t *testing.T) {
	adminInfo, err := tester.LoginAdmin()

	assert.Nil(t, err)

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

		n := schema.Notification{}
		assert.Nil(t, r.Decode(&n))
		assert.Equal(t, newTittle, n.Title)
		assert.Equal(t, newContent, n.Content)
		assert.Equal(t, newNote, *n.Note)
	}
}
