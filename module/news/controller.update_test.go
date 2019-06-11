// Copyright 2019 Axetroy. All rights reserved. MIT license.
package news_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_schema"
	"github.com/axetroy/go-server/module/news"
	"github.com/axetroy/go-server/module/news/news_model"
	"github.com/axetroy/go-server/module/news/news_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUpdate(t *testing.T) {
	// 更新成功
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

		adminInfo := admin_schema.AdminProfileWithToken{}

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
			newsType = news_model.NewsType_News
		)

		r := news.Create(schema.Context{
			Uid: adminUid,
		}, news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    newsType,
			Tags:    []string{},
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := news_schema.News{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		newsId = n.Id

		defer news.DeleteNewsById(n.Id)
	}

	// 3. 更新这篇新闻公告
	{
		var (
			newTittle  = "new title"
			newContent = "new content"
			newType    = news_model.NewsType_Announcement
			newTags    = []string{newTittle}
		)

		r := news.Update(schema.Context{
			Uid: adminUid,
		}, newsId, news.UpdateParams{
			Title:   &newTittle,
			Content: &newContent,
			Type:    &newType,
			Tags:    &newTags,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		newsInfo := news_schema.News{}

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
		r2 := news.Update(schema.Context{
			Uid: adminUid,
		}, newsId, news.UpdateParams{
			Title: &newTittle2,
		})

		assert.Equal(t, schema.StatusSuccess, r2.Status)
		assert.Equal(t, "", r2.Message)

		newsInfo2 := news_schema.News{}

		assert.Nil(t, tester.Decode(r2.Data, &newsInfo2))

		assert.Equal(t, newTittle2, newsInfo2.Title)
		assert.Equal(t, newContent, newsInfo2.Content)
		assert.Equal(t, newType, newsInfo2.Type)
		assert.Equal(t, newTags, newsInfo2.Tags)
	}
}

func TestUpdateRouter(t *testing.T) {
	var (
		newsId string
	)

	adminInfo, _ := tester.LoginAdmin()

	// 先创建一篇新闻作为测试
	{
		var (
			title    = "test"
			content  = "test"
			newsType = news_model.NewsType_News
		)

		r := news.Create(schema.Context{
			Uid: adminInfo.Id,
		}, news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    newsType,
			Tags:    []string{},
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := news_schema.News{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		newsId = n.Id

		defer news.DeleteNewsById(n.Id)
	}

	var (
		newTitle   = "new title"
		newContent = "new content"
	)

	// 更新
	{

		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		body, _ := json.Marshal(&news.UpdateParams{
			Title:   &newTitle,
			Content: &newContent,
		})

		r := tester.HttpAdmin.Put("/v1/news/n/"+newsId, body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		n := news_schema.News{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		assert.Equal(t, newTitle, n.Title)
		assert.Equal(t, newContent, n.Content)
	}

	// 获取详情查看是否更改成功
	{
		res := news.GetNews(newsId)

		n := news_schema.News{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		assert.Equal(t, newTitle, n.Title)
		assert.Equal(t, newContent, n.Content)
	}
}
