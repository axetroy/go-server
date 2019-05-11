package message_test

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/message"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
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

func TestGetMessage(t *testing.T) {
	// 获取一篇存在的消息公告
	{
		var (
			adminUid  string
			messageId string
		)
		// 1. 先登陆获取管理员的Token
		{
			r := admin.Login(admin.SignInParams{
				Username: "admin",
				Password: "admin",
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			adminInfo := schema.AdminProfileWithToken{}

			assert.Nil(t, tester.Decode(r.Data, &adminInfo))

			if c, er := util.ParseToken(util.TokenPrefix+" "+adminInfo.Token, true); er != nil {
				t.Error(er)
			} else {
				adminUid = c.Uid
			}
		}

		// 2. 先创建一篇消息作为测试
		{
			var (
				title   = "test"
				content = "test"
			)

			r := message.Create(controller.Context{
				Uid: adminUid,
			}, message.CreateMessageParams{
				Uid:     adminUid,
				Title:   title,
				Content: content,
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			n := schema.Message{}

			assert.Nil(t, tester.Decode(r.Data, &n))

			messageId = n.Id

			defer message.DeleteMessageById(n.Id)
		}

		// 3. 获取文章公告
		{
			r := message.Get(controller.Context{
				Uid: adminUid,
			}, messageId)

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			messageInfo := r.Data.(schema.Message)

			assert.Equal(t, "test", messageInfo.Title)
			assert.Equal(t, "test", messageInfo.Content)
		}
	}
}
