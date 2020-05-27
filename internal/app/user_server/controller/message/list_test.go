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

func TestGetMessageListByUser(t *testing.T) {

	{
		var (
			data = make([]schema.Message, 0)
		)
		query := schema.Query{
			Limit: 20,
		}
		r := message.GetMessageListByUser(helper.Context{
			Uid: "123123",
		}, message.Query{
			Query: query,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, r.Decode(&data))
		assert.Equal(t, query.Limit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.Equal(t, 0, r.Meta.Num)
		assert.Equal(t, int64(0), r.Meta.Total)
	}

	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

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

		defer message.DeleteMessageById(n.Id)
	}

	// 3. 获取列表
	{
		data := make([]schema.Message, 0)

		query := schema.Query{
			Limit: 20,
		}
		r := message.GetMessageListByUser(helper.Context{
			Uid: userInfo.Id,
		}, message.Query{
			Query: query,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, r.Decode(&data))

		assert.Equal(t, query.Limit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.Equal(t, 1, r.Meta.Num)
		assert.Equal(t, int64(1), r.Meta.Total)

		assert.Len(t, data, 1)
	}
}

func TestGetMessageListByUserRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

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

		//defer message.DeleteMessageById(n.Id)
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	{
		r := tester.HttpUser.Get("/v1/message", nil, &header)

		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)

		assert.Equal(t, "", res.Message)

		messages := make([]schema.Message, 0)

		assert.Nil(t, res.Decode(&messages))

		assert.True(t, len(messages) > 0)

		for _, b := range messages {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}
