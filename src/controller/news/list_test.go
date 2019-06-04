// Copyright 2019 Axetroy. All rights reserved. MIT license.
package news_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/news"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	// 现在没有任何文章，获取到的应该是0个长度的
	{
		var (
			data = make([]model.News, 0)
		)
		query := schema.Query{
			Limit: 20,
		}
		r := news.GetList(news.Query{
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

	adminInfo, _ := tester.LoginAdmin()

	// 2. 先创建一篇新闻作为测试
	{
		var (
			title    = "test"
			content  = "test"
			newsType = model.NewsType_News
		)

		r := news.Create(controller.Context{
			Uid: adminInfo.Id,
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

		defer news.DeleteNewsById(n.Id)
	}

	// 3. 获取列表
	{
		var (
			data = make([]model.News, 0)
		)
		query := schema.Query{
			Limit: 20,
		}
		r := news.GetList(news.Query{
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

func TestGetListRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	{
		var (
			title    = "test"
			content  = "test"
			newsType = model.NewsType_News
		)

		r := news.Create(controller.Context{
			Uid: adminInfo.Id,
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

		defer news.DeleteNewsById(n.Id)
	}

	{
		r := tester.HttpUser.Get("/v1/news", nil, nil)

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		banners := make([]schema.News, 0)

		assert.Nil(t, tester.Decode(res.Data, &banners))

		for _, b := range banners {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, "string", b.Author)
			assert.IsType(t, model.NewsType_Announcement, b.Type)
			assert.IsType(t, []string{""}, b.Tags)
			assert.IsType(t, model.NewsStatusActive, b.Status)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}
