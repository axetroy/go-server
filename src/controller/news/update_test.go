package news_test

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/news"
	"github.com/axetroy/go-server/src/model"
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

func TestUpdate(t *testing.T) {
	// 更新成功
	{
		// 1. 先登陆获取管理员的Token
		var (
			adminUid string
			newsId   string
		)
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

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			n := schema.News{}

			assert.Nil(t, tester.Decode(r.Data, &n))

			newsId = n.Id

			defer func() {
				news.DeleteNewsById(n.Id)
			}()
		}

		// 3. 更新这篇新闻公告
		{
			var (
				newTittle  = "new title"
				newContent = "new content"
				newType    = model.NewsType_Announcement
				newTags    = []string{newTittle}
			)

			r := news.Update(controller.Context{
				Uid: adminUid,
			}, newsId, news.UpdateParams{
				Tittle:  &newTittle,
				Content: &newContent,
				Type:    &newType,
				Tags:    &newTags,
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			newsInfo := schema.News{}

			assert.Nil(t, tester.Decode(r.Data, &newsInfo))

			assert.Equal(t, newTittle, newsInfo.Title)
			assert.Equal(t, newContent, newsInfo.Content)
			assert.Equal(t, newType, newsInfo.Type)
			assert.Equal(t, newTags, newsInfo.Tags)

			var (
				newTittle2 = "new title 2"
			)

			// 只更新部分字段
			// 其余字段应该保持不变
			r2 := news.Update(controller.Context{
				Uid: adminUid,
			}, newsId, news.UpdateParams{
				Tittle: &newTittle2,
			})

			assert.Equal(t, schema.StatusSuccess, r2.Status)
			assert.Equal(t, "", r2.Message)

			newsInfo2 := schema.News{}

			assert.Nil(t, tester.Decode(r2.Data, &newsInfo2))

			assert.Equal(t, newTittle2, newsInfo2.Title)
			assert.Equal(t, newContent, newsInfo2.Content)
			assert.Equal(t, newType, newsInfo2.Type)
			assert.Equal(t, newTags, newsInfo2.Tags)
		}
	}
}
