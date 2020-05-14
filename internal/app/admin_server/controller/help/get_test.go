// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package help_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/help"
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

func TestGetHelp(t *testing.T) {
	{
		r := help.GetHelp("123123")

		assert.Equal(t, exception.NoData.Code(), r.Status)
		assert.Equal(t, exception.NoData.Error(), r.Message)
	}

	{
		var (
			helpId  string
			title   = "test"
			content = "test"
		)

		adminInfo, _ := tester.LoginAdmin()

		// 2. 创建作为测试
		{

			r := help.Create(helper.Context{
				Uid: adminInfo.Id,
			}, help.CreateParams{
				Title:   title,
				Content: content,
				Tags:    []string{},
				Status:  model.HelpStatusActive,
				Type:    model.HelpTypeArticle,
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			n := schema.Help{}

			assert.Nil(t, r.Decode(&n))

			helpId = n.Id

			defer help.DeleteHelpById(n.Id)
		}

		// 3. 获取文章公告
		{
			r := help.GetHelp(helpId)

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			helpInfo := schema.Help{}

			assert.Nil(t, r.Decode(&helpInfo))

			assert.Equal(t, title, helpInfo.Title)
			assert.Equal(t, content, helpInfo.Content)
		}
	}
}

func TestGetHelpRouter(t *testing.T) {
	var (
		helpId  = ""
		title   = "test title"
		content = "test content"
		tags    = []string{"test"}
	)

	adminInfo, _ := tester.LoginAdmin()

	// 创建一个 help
	{

		r := help.Create(helper.Context{
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

		assert.Nil(t, r.Decode(&n))

		defer help.DeleteHelpById(n.Id)

		helpId = n.Id
		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
		assert.Equal(t, tags, n.Tags)
		assert.Equal(t, model.HelpStatusActive, n.Status)
		assert.Equal(t, model.HelpTypeArticle, n.Type)
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	// 获取详情
	{

		r := tester.HttpAdmin.Get("/v1/help/h/"+helpId, nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		n := schema.Help{}

		assert.Nil(t, res.Decode(&n))

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
		assert.Equal(t, tags, n.Tags)
		assert.Equal(t, model.HelpTypeArticle, n.Type)
		assert.Equal(t, model.HelpStatusActive, n.Status)
	}
}
