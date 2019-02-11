package message_test

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/message"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestMarkRead(t *testing.T) {
	var (
		adminUid string
	)
	// 先登陆获取管理员的Token
	{
		// 登陆超级管理员-成功

		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "admin",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := schema.AdminProfileWithToken{}

		if err := tester.Decode(r.Data, &adminInfo); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, "admin", adminInfo.Username)
		assert.True(t, len(adminInfo.Token) > 0)

		if c, er := util.ParseToken(util.TokenPrefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			adminUid = c.Uid
		}
	}

	context := controller.Context{
		Uid: adminUid,
	}

	var testMessage schema.Message

	var testUser schema.Profile

	{
		// 创建一个测试用户
		// 1。 创建测试账号
		rand.Seed(111)
		username := "test-TestMarkRead"
		password := "123123"

		r := auth.SignUp(auth.SignUpParams{
			Username: &username,
			Password: password,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = schema.Profile{}

		if err := tester.Decode(r.Data, &testUser); err != nil {
			t.Error(err)
			return
		}

		defer auth.DeleteUserByUserName(username)
	}

	// 创建一篇个人消息
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := message.Create(context, message.CreateMessageParams{
			Uid:     testUser.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testMessage = schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &testMessage))

		defer message.DeleteMessageById(testMessage.Id)

		assert.Equal(t, title, testMessage.Tittle)
		assert.Equal(t, content, testMessage.Content)
	}

	{
		// 用测试用户标记为已读
		r := message.MarkRead(controller.Context{
			Uid: testUser.Id,
		}, testMessage.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
	}

	{
		// 用测试者的账号获取详情
		r := message.Get(controller.Context{
			Uid: testUser.Id,
		}, testMessage.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := r.Data.(schema.Message)

		assert.Equal(t, testMessage.Id, n.Id)
		assert.Equal(t, testMessage.Tittle, n.Tittle)
		assert.Equal(t, testMessage.Content, n.Content)
		assert.Equal(t, true, n.Read)
		assert.IsType(t, "", *n.ReadAt)
	}
}

func TestReadRouter(t *testing.T) {
	// TODO: 完善HTTP的测试用例
	t.Skip()
}
