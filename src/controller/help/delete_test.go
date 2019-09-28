// Copyright 2019 Axetroy. All rights reserved. MIT license.
package help_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/help"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDelete(t *testing.T) {
	var (
		helpId  = ""
		title   = "test title"
		content = "test content"
		tags    = []string{"test"}
	)

	adminInfo, _ := tester.LoginAdmin()

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

		helpId = n.Id
		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
		assert.Equal(t, tags, n.Tags)
		assert.Equal(t, model.HelpStatusActive, n.Status)
		assert.Equal(t, model.HelpTypeArticle, n.Type)
	}

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	// 删除这个刚添加的地址
	{
		r := help.Delete(context, helpId)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		helpInfo := schema.Help{}

		assert.Nil(t, tester.Decode(r.Data, &helpInfo))

		assert.Equal(t, title, helpInfo.Title)
		assert.Equal(t, content, helpInfo.Content)
		assert.Equal(t, tags, helpInfo.Tags)

		if err := database.Db.First(&model.Banner{
			Id: helpId,
		}).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				assert.Fail(t, "数据被删除，应该不能再找到")
			}
		} else {
			assert.Fail(t, "数据被删除，应该不能再找到")
		}
	}

}

func TestDeleteRouter(t *testing.T) {
	var (
		helpId  = ""
		title   = "test title"
		content = "test content"
		tags    = []string{"test"}
	)

	adminInfo, _ := tester.LoginAdmin()

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

	// 删除这条地址
	{
		r := tester.HttpAdmin.Delete("/v1/help/h/"+helpId, nil, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		helpInfo := schema.Help{}

		assert.Nil(t, tester.Decode(res.Data, &helpInfo))

		assert.Equal(t, title, helpInfo.Title)
		assert.Equal(t, content, helpInfo.Content)
		assert.Equal(t, tags, helpInfo.Tags)

	}

}
