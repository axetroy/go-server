// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package news_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/news"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	var (
		adminInfo schema.AdminProfileWithToken
		userInfo  schema.ProfileWithToken
		err       error
	)
	adminInfo, err = tester.LoginAdmin()

	if !assert.Nil(t, err) {
		return
	}

	userInfo, err = tester.CreateUser()

	if !assert.Nil(t, err) {
		return
	}

	defer tester.DeleteUserByUserName(userInfo.Username)

	// 创建一个公告
	{
		var (
			title   = "test"
			content = "test"
		)

		r := news.Create(helper.Context{
			Uid: adminInfo.Id,
		}, news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    model.NewsTypeNews,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := model.News{}

		assert.Nil(t, r.Decode(&n))

		defer news.DeleteNewsById(n.Id)

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
	}

	// 非管理员的uid去创建，应该报错
	{
		var (
			title    = "test"
			content  = "test"
			newsType = model.NewsTypeNews
		)

		r := news.Create(helper.Context{
			Uid: userInfo.Id,
		}, news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    newsType,
			Tags:    []string{},
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.AdminNotExist.Error(), r.Message)
	}
}

func TestCreateRouter(t *testing.T) {
	var (
		adminInfo schema.AdminProfileWithToken
		err       error
	)

	adminInfo, err = tester.LoginAdmin()

	if !assert.Nil(t, err) {
		return
	}

	{
		var (
			title   = "test"
			content = "test"
		)

		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		body, _ := json.Marshal(&news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    model.NewsTypeNews,
		})

		r := tester.HttpAdmin.Post("/v1/news", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		n := schema.News{}

		assert.Nil(t, res.Decode(&n))

		defer news.DeleteNewsById(n.Id)

		assert.Equal(t, adminInfo.Id, n.Author)
		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
	}
}
