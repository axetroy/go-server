package message_test

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/message"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	// 现在没有任何文章，获取到的应该是0个长度的
	{
		var (
			data = make([]model.Message, 0)
		)
		query := schema.Query{
			Limit: 20,
		}
		r := message.GetList(message.Query{
			Query: query,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &data))
		assert.Equal(t, query.Limit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.Equal(t, 0, r.Meta.Num)
		assert.Equal(t, int64(0), r.Meta.Total)
	}

	{
		var (
			adminUid string
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

		// 2. 先创建一篇新闻作为测试
		{
			var (
				title   = "test"
				content = "test"
			)

			r := message.Create(controller.Context{
				Uid: adminUid,
			}, message.CreateMessageParams{
				Title:   title,
				Content: content,
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			n := schema.Message{}

			assert.Nil(t, tester.Decode(r.Data, &n))

			defer message.DeleteMessageById(n.Id)
		}

		// 3. 获取列表
		{
			var (
				data = make([]model.Message, 0)
			)
			query := schema.Query{
				Limit: 20,
			}
			r := message.GetList(message.Query{
				Query: query,
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			assert.Nil(t, tester.Decode(r.Data, &data))
			assert.Equal(t, query.Limit, r.Meta.Limit)
			assert.Equal(t, schema.DefaultPage, r.Meta.Page)
			assert.Equal(t, 1, r.Meta.Num)
			assert.Equal(t, int64(1), r.Meta.Total)

			assert.Len(t, data, 1)
		}
	}
}
