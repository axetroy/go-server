package news_test

import (
	"github.com/axetroy/go-server/controller"
	"github.com/axetroy/go-server/controller/admin"
	"github.com/axetroy/go-server/controller/news"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/go-server/token"
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

func TestGetNews(t *testing.T) {
	// 获取一篇不存在的新闻公告
	{
		r := news.GetNews("123123")

		assert.Equal(t, response.StatusFail, r.Status)
		assert.Equal(t, exception.NewsNotExist.Error(), r.Message)
	}

	// 获取一篇存在的新闻公告
	{
		var (
			adminUid string
			newsId   string
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

			r := news.Create(controller.Context{
				Uid: adminUid,
			}, news.CreateNewParams{
				Title:   title,
				Content: content,
				Type:    newsType,
				Tags:    []string{},
			})

			assert.Equal(t, response.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			n := news.News{}

			assert.Nil(t, tester.Decode(r.Data, &n))

			newsId = n.Id

			defer func() {
				news.DeleteNewsById(n.Id)
			}()
		}

		// 3. 获取文章公告
		{
			r := news.GetNews(newsId)

			assert.Equal(t, response.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			newsInfo := news.News{}

			assert.Nil(t, tester.Decode(r.Data, &newsInfo))

			assert.Equal(t, "test", newsInfo.Tittle)
			assert.Equal(t, "test", newsInfo.Content)
			assert.Equal(t, model.NewsType_News, newsInfo.Type)
		}
	}
}
