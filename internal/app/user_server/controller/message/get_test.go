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

func TestGetMessage(t *testing.T) {
	var (
		messageId string
	)

	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	// 2. 创建
	{
		var (
			title   = "test"
			content = "test"
		)

		r := messageAdmin.Create(helper.Context{
			Uid: adminInfo.Id,
		}, messageAdmin.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Message{}

		assert.Nil(t, r.Decode(&n))

		messageId = n.Id

		defer message.DeleteMessageById(n.Id)
	}

	// 3. 获取
	{
		r := message.Get(helper.Context{
			Uid: userInfo.Id,
		}, messageId)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		messageInfo := r.Data.(schema.Message)

		assert.Equal(t, "test", messageInfo.Title)
		assert.Equal(t, "test", messageInfo.Content)
	}

}

func TestGetRouter(t *testing.T) {
	var (
		messageId string
	)

	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	// 2. 先创建一篇消息作为测试
	{
		var (
			title   = "test"
			content = "test"
		)

		r := messageAdmin.Create(helper.Context{
			Uid: adminInfo.Id,
		}, messageAdmin.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Message{}

		assert.Nil(t, r.Decode(&n))

		messageId = n.Id

		defer message.DeleteMessageById(n.Id)
	}

	// 用户接口获取
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + userInfo.Token,
		}

		r := tester.HttpUser.Get("/v1/message/"+messageId, nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		n := schema.Message{}

		assert.Nil(t, res.Decode(&n))

		assert.Equal(t, "test", n.Title)
		assert.Equal(t, "test", n.Content)
		assert.IsType(t, "string", n.CreatedAt)
		assert.IsType(t, "string", n.UpdatedAt)
		assert.False(t, n.Read)
		assert.Nil(t, n.ReadAt)
	}

}
