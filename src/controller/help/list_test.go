// Copyright 2019 Axetroy. All rights reserved. MIT license.
package help_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/help"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	context := controller.Context{
		Uid: adminInfo.Id,
	}
	var (
		title   = "test title"
		content = "test content"
		tags    = []string{"test"}
	)

	// 创建一个 help
	{

		r := help.Create(controller.Context{
			Uid: adminInfo.Id,
		}, help.CreateParams{
			Title:   title,
			Content: content,
			Tags:    tags,
			Status:  model.HelpStatusActive,
			Type:    model.HelpTypeArticle,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Help{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer help.DeleteHelpById(n.Id)

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
		assert.Equal(t, tags, n.Tags)
		assert.Equal(t, model.HelpStatusActive, n.Status)
		assert.Equal(t, model.HelpTypeArticle, n.Type)
	}

	// 获取列表
	{
		r := help.GetHelpList(context, help.Query{})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		helps := make([]schema.Help, 0)

		assert.Nil(t, tester.Decode(r.Data, &helps))

		assert.Equal(t, schema.DefaultLimit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.IsType(t, 1, r.Meta.Num)
		assert.IsType(t, int64(1), r.Meta.Total)

		assert.True(t, len(helps) >= 1)

		for _, b := range helps {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, model.HelpTypeArticle, b.Type)
			assert.IsType(t, model.HelpStatusActive, b.Status)
			assert.IsType(t, []string{}, b.Tags)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}

func TestGetListRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	{
		var (
			title   = "test title"
			content = "test content"
			tags    = []string{"test"}
		)

		// 创建一个 help
		r := help.Create(controller.Context{
			Uid: adminInfo.Id,
		}, help.CreateParams{
			Title:   title,
			Content: content,
			Tags:    tags,
			Status:  model.HelpStatusActive,
			Type:    model.HelpTypeArticle,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Help{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer help.DeleteHelpById(n.Id)

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
		assert.Equal(t, tags, n.Tags)
		assert.Equal(t, model.HelpStatusActive, n.Status)
		assert.Equal(t, model.HelpTypeArticle, n.Type)
	}

	{
		r := tester.HttpAdmin.Get("/v1/help", nil, &header)

		res := schema.List{}

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		helps := make([]schema.Help, 0)

		assert.Nil(t, tester.Decode(res.Data, &helps))

		assert.Equal(t, schema.DefaultLimit, res.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, res.Meta.Page)
		assert.IsType(t, 1, res.Meta.Num)
		assert.IsType(t, int64(1), res.Meta.Total)

		for _, b := range helps {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, model.HelpTypeArticle, b.Type)
			assert.IsType(t, model.HelpStatusActive, b.Status)
			assert.IsType(t, []string{}, b.Tags)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}
