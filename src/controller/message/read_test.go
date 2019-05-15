package message_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/message"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMarkRead(t *testing.T) {
	var (
		testMessage schema.Message
	)

	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	// 创建一篇个人消息
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

		testMessage = schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &testMessage))

		defer message.DeleteMessageById(testMessage.Id)

		assert.Equal(t, title, testMessage.Title)
		assert.Equal(t, content, testMessage.Content)
	}

	{
		// 用测试用户标记为已读
		r := message.MarkRead(controller.Context{
			Uid: userInfo.Id,
		}, testMessage.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
	}

	{
		// 用测试者的账号获取详情
		r := message.Get(controller.Context{
			Uid: userInfo.Id,
		}, testMessage.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := r.Data.(schema.Message)

		assert.Equal(t, testMessage.Id, n.Id)
		assert.Equal(t, testMessage.Title, n.Title)
		assert.Equal(t, testMessage.Content, n.Content)
		assert.Equal(t, true, n.Read)
		assert.IsType(t, "", *n.ReadAt)
	}
}

func TestReadRouter(t *testing.T) {
	var (
		testMessage schema.Message
	)

	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	// 创建一篇个人消息
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

		testMessage = schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &testMessage))

		defer message.DeleteMessageById(testMessage.Id)

		assert.Equal(t, title, testMessage.Title)
		assert.Equal(t, content, testMessage.Content)
	}

	// 用户标记为已读
	{

		header := mocker.Header{
			"Authorization": util.TokenPrefix + " " + userInfo.Token,
		}

		r := tester.HttpUser.Put("/v1/message/m/"+testMessage.Id+"/read", nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		n := schema.Message{}

		assert.Nil(t, tester.Decode(res.Data, &n))
	}
}
