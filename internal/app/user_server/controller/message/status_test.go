// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message_test

import (
	"encoding/json"
	messageAdmin "github.com/axetroy/go-server/internal/app/admin_server/controller/message"
	"github.com/axetroy/go-server/internal/app/user_server/controller/message"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetStatusRouter(t *testing.T) {
	var (
		testMessage schema.Message
	)

	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	context := helper.Context{
		Uid: adminInfo.Id,
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	// 创建一篇个人消息
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := messageAdmin.Create(context, messageAdmin.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testMessage = schema.Message{}

		assert.Nil(t, r.Decode(&testMessage))

		defer message.DeleteMessageById(testMessage.Id)

		assert.Equal(t, title, testMessage.Title)
		assert.Equal(t, content, testMessage.Content)
	}

	// 获取未读消息
	{
		r := tester.HttpUser.Get("/v1/message/status", nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		n := schema.MessageStatus{}

		assert.Nil(t, res.Decode(&n))
		assert.Equal(t, n.Unread, int64(1))
	}

	// 用户标记为已读
	{
		r := tester.HttpUser.Put("/v1/message/"+testMessage.Id+"/read", nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		n := schema.Message{}

		assert.Nil(t, res.Decode(&n))
		assert.Equal(t, true, n.Read)
		assert.NotNil(t, n.ReadAt)
	}

	// 获取未读消息
	{
		r := tester.HttpUser.Get("/v1/message/status", nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		n := schema.MessageStatus{}

		assert.Nil(t, res.Decode(&n))
		assert.Equal(t, n.Unread, int64(0))
	}
}
