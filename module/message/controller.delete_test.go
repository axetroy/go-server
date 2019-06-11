// Copyright 2019 Axetroy. All rights reserved. MIT license.
package message_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/message"
	"github.com/axetroy/go-server/module/message/message_model"
	"github.com/axetroy/go-server/module/message/message_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDeleteByAdmin(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := schema.Context{
		Uid: adminInfo.Id,
	}

	var testMessage message_schema.Message

	// 创建一条个人信息
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := message.Create(context, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testMessage = message_schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &testMessage))

		defer message.DeleteMessageById(testMessage.Id)

		assert.Equal(t, title, testMessage.Title)
		assert.Equal(t, content, testMessage.Content)
	}

	// 获取通知
	{
		n := message_model.Message{
			Id: testMessage.Id,
		}

		assert.Nil(t, database.Db.Model(&n).Where(&n).First(&n).Error)
	}

	// 删除通知
	{
		r := message.DeleteByAdmin(context, testMessage.Id)

		assert.Equal(t, "", r.Message)
		assert.Equal(t, schema.StatusSuccess, r.Status)
	}

	// 再次获取通知，这时候通知应该已经被删除了
	{
		n := message_model.Message{
			Id: testMessage.Id,
		}

		if err := database.Db.First(&n).Error; err != nil {
			assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
		} else {
			assert.Fail(t, "通知应该已被删除")
		}
	}
}

func TestDeleteByAdminRouter(t *testing.T) {
	var (
		messageId string
	)

	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	// 创建一条系统通知
	{
		var (
			title   = "test title"
			content = "test content"
		)

		body, _ := json.Marshal(&message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		r := tester.HttpAdmin.Post("/v1/message", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		messageInfo := message_schema.Message{}

		assert.Nil(t, tester.Decode(res.Data, &messageInfo))

		messageId = messageInfo.Id

		defer message.DeleteMessageById(messageInfo.Id)

		assert.Equal(t, title, messageInfo.Title)
		assert.Equal(t, content, messageInfo.Content)
	}

	// 删除这条通知
	{
		r := tester.HttpAdmin.Delete("/v1/message/m/"+messageId, nil, &header)

		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		// 再查找这条记录，应该是空的

		n := message_model.Message{Id: messageId}

		err := database.Db.First(&n).Error

		assert.NotNil(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
	}
}

func TestDeleteByUser(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := schema.Context{
		Uid: adminInfo.Id,
	}

	var testMessage message_schema.Message

	// 创建一条个人信息
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := message.Create(context, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testMessage = message_schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &testMessage))

		defer message.DeleteMessageById(testMessage.Id)

		assert.Equal(t, title, testMessage.Title)
		assert.Equal(t, content, testMessage.Content)
	}

	// 获取消息
	{
		n := message_model.Message{
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
		n := message_model.Message{
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
		messageInfo = message_schema.Message{}
	)

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

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

		r := message.Create(schema.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		messageInfo = message_schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &messageInfo))

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

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
			return
		}

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		// 再查找这条记录，应该是空的

		n := message_model.Message{
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
