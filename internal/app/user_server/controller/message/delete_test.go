// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message_test

import (
	"encoding/json"
	messageAdmin "github.com/axetroy/go-server/internal/app/admin_server/controller/message"
	"github.com/axetroy/go-server/internal/app/user_server/controller/message"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDeleteByUser(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	context := helper.Context{
		Uid: adminInfo.Id,
	}

	var testMessage schema.Message

	// 创建一条个人信息
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

	// 获取消息
	{
		n := model.Message{
			Id: testMessage.Id,
		}

		assert.Nil(t, database.Db.Model(&n).Where(&n).First(&n).Error)
	}

	// 删除消息
	{
		r := message.DeleteByUser(context, testMessage.Id)

		assert.Equal(t, "", r.Message)
		assert.Equal(t, schema.StatusSuccess, r.Status)
	}

	// 再次获取通知，这时候通知应该已经被删除了
	{
		n := model.Message{
			Id: testMessage.Id,
		}

		if err := database.Db.First(&n).Error; err != nil {
			assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
		} else {
			assert.Fail(t, "通知应该已被删除")
		}
	}
}

func TestDeleteByUserRouter(t *testing.T) {
	var (
		messageInfo = schema.Message{}
	)

	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	adminInfo, _ := tester.LoginAdmin()

	userHeader := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	// 创建一条个人信息
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
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

		messageInfo = schema.Message{}

		assert.Nil(t, r.Decode(&messageInfo))

		defer message.DeleteMessageById(messageInfo.Id)

		assert.Equal(t, title, messageInfo.Title)
		assert.Equal(t, content, messageInfo.Content)
	}

	// 删除这条通知
	{
		r := tester.HttpUser.Delete("/v1/message/m/"+messageInfo.Id, nil, &userHeader)

		res := schema.Response{}

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		// 再查找这条记录，应该是空的

		n := model.Message{
			Id: messageInfo.Id,
		}

		err := database.Db.First(&n).Error

		if !assert.NotNil(t, err) {
			return
		}
		if !assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error()) {
			return
		}
	}
}
