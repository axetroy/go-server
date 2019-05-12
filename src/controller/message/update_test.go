package message_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/message"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func init() {
	// 确保超级管理员存在
	admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "admin",
		Password: "admin",
		Name:     "admin",
	}, true)
}

func TestUpdate(t *testing.T) {
	var (
		messageInfo = schema.Message{}
	)

	adminInfo, _ := tester.LoginAdmin()

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	// 创建一个消息
	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(controller.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := model.Message{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer message.DeleteMessageById(n.Id)

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
	}

	// 更新这个刚添加的消息
	{

		var (
			newTitle   = "new title"
			newContent = "new content"
		)

		r := message.Update(context, messageInfo.Id, message.UpdateParams{
			Title:   &newTitle,
			Content: &newContent,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &messageInfo))

		assert.Equal(t, newTitle, messageInfo.Title)
		assert.Equal(t, newContent, messageInfo.Content)
	}
}

func TestUpdateRouter(t *testing.T) {
	var (
		messageInfo = schema.Message{}
	)

	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": util.TokenPrefix + " " + adminInfo.Token,
	}

	// 创建一个消息
	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(controller.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &messageInfo))

		defer message.DeleteMessageById(messageInfo.Id)

		assert.Equal(t, title, messageInfo.Title)
		assert.Equal(t, content, messageInfo.Content)
	}

	// 修改这条 banner
	{

		var (
			newTitle   = "new title"
			newContent = "new content"
		)

		body, _ := json.Marshal(&message.UpdateParams{
			Title:   &newTitle,
			Content: &newContent,
		})

		r := tester.HttpAdmin.Put("/v1/message/update/"+messageInfo.Id, body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

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

		assert.Nil(t, tester.Decode(res.Data, &messageInfo))

		assert.Equal(t, newTitle, messageInfo.Title)
		assert.Equal(t, newContent, messageInfo.Content)

	}

}
