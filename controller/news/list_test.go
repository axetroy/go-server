package news_test

import (
	"github.com/axetroy/go-server/controller/admin"
	"github.com/axetroy/go-server/controller/news"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/request"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/go-server/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	// 现在没有任何文章，获取到的应该是0个长度的
	{
		var (
			data = make([]model.News, 0)
		)
		query := request.Query{
			Limit: 20,
		}
		r := news.GetList(news.Query{
			Query: query,
		})

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &data))
		assert.Equal(t, query.Limit, r.Meta.Limit)
		assert.Equal(t, request.DefaultPage, r.Meta.Page)
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

			assert.Equal(t, response.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			adminInfo := admin.SignInResponse{}

			assert.Nil(t, tester.Decode(r.Data, &adminInfo))

			if c, er := token.Parse(token.Prefix+" "+adminInfo.Token, true); er != nil {
				t.Error(er)
			} else {
				adminUid = c.Uid
			}
		}

		// 2. 先创建一篇新闻作为测试
		{
			var (
				title    = "test"
				content  = "test"
				newsType = model.NewsType_News
			)

			r := news.Create(adminUid, news.CreateNewParams{
				Title:   title,
				Content: content,
				Type:    newsType,
				Tags:    []string{},
			})

			assert.Equal(t, response.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			n := news.News{}

			assert.Nil(t, tester.Decode(r.Data, &n))

			defer func() {
				news.DeleteNewsById(n.Id)
			}()
		}

		// 3. 获取列表
		{
			var (
				data = make([]model.News, 0)
			)
			query := request.Query{
				Limit: 20,
			}
			r := news.GetList(news.Query{
				Query: query,
			})

			assert.Equal(t, response.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			assert.Nil(t, tester.Decode(r.Data, &data))
			assert.Equal(t, query.Limit, r.Meta.Limit)
			assert.Equal(t, request.DefaultPage, r.Meta.Page)
			assert.Equal(t, 1, r.Meta.Num)
			assert.Equal(t, int64(1), r.Meta.Total)

			assert.Len(t, data, 1)
		}
	}
}
