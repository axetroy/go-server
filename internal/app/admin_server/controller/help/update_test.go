// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package help_test

import (
	"encoding/json"
	help2 "github.com/axetroy/go-server/internal/app/admin_server/controller/help"
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

func TestUpdate(t *testing.T) {
	var (
		helpId  = ""
		title   = "test title"
		content = "test content"
		tags    = []string{"test"}
	)

	adminInfo, _ := tester.LoginAdmin()

	// 创建一个 help
	{

		r := help2.Create(helper.Context{
			Uid: adminInfo.Id,
		}, help2.CreateParams{
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

		defer help2.DeleteHelpById(n.Id)

		helpId = n.Id
		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
		assert.Equal(t, tags, n.Tags)
		assert.Equal(t, model.HelpStatusActive, n.Status)
		assert.Equal(t, model.HelpTypeArticle, n.Type)
	}

	context := helper.Context{
		Uid: adminInfo.Id,
	}

	// 更新
	{
		var (
			newTitle = "new address"
		)

		r := help2.Update(context, helpId, help2.UpdateParams{
			Title: &newTitle,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		helpInfo := schema.Help{}

		assert.Nil(t, r.Decode(&helpInfo))

		assert.Equal(t, newTitle, helpInfo.Title)
	}

	// 再次更新
	{
		var (
			newTitle   = "new address"
			newContent = "http://test.com"
		)

		r := help2.Update(context, helpId, help2.UpdateParams{
			Title:   &newTitle,
			Content: &newContent,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		helpInfo := schema.Help{}

		assert.Nil(t, r.Decode(&helpInfo))

		assert.Equal(t, newTitle, helpInfo.Title)
		assert.Equal(t, newContent, helpInfo.Content)
	}
}

func TestUpdateRouter(t *testing.T) {
	var (
		helpId  = ""
		title   = "test title"
		content = "test content"
		tags    = []string{"test"}
	)

	adminInfo, _ := tester.LoginAdmin()

	// 创建一个 help
	{

		r := help2.Create(helper.Context{
			Uid: adminInfo.Id,
		}, help2.CreateParams{
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

		defer help2.DeleteHelpById(n.Id)

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

	// 修改
	{

		var (
			newTitle   = "new.png"
			newContent = "13333333333"
		)

		body, _ := json.Marshal(&help2.UpdateParams{
			Title:   &newTitle,
			Content: &newContent,
		})

		r := tester.HttpAdmin.Put("/v1/help/"+helpId, body, &header)

		assert.Equal(t, http.StatusOK, r.Code)

		res := schema.Response{}

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, "", res.Message)
		assert.Equal(t, schema.StatusSuccess, res.Status)

		helpInfo := schema.Help{}

		assert.Nil(t, res.Decode(&helpInfo))

		assert.Equal(t, newTitle, helpInfo.Title)
		assert.Equal(t, newContent, helpInfo.Content)

	}

}
