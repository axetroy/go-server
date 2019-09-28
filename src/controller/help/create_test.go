package help_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/help"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	// 创建一个 help
	{
		var (
			title   = "test title"
			content = "test content"
			tags    = []string{"test"}
		)

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

	// 非管理员的uid去创建，应该报错
	{

		userInfo, _ := tester.CreateUser()

		defer auth.DeleteUserByUserName(userInfo.Username)

		var (
			title   = "test title"
			content = "test content"
			tags    = []string{"test"}
		)

		r := help.Create(controller.Context{
			Uid: userInfo.Id,
		}, help.CreateParams{
			Title:   title,
			Content: content,
			Tags:    tags,
			Status:  model.HelpStatusActive,
			Type:    model.HelpTypeArticle,
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.AdminNotExist.Error(), r.Message)
	}
}

func TestCreateRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	// 创建 help
	{
		var (
			tittle  = "tittle"
			content = "content"
		)

		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		body, _ := json.Marshal(&help.CreateParams{
			Title:   tittle,
			Content: content,
			Status:  model.HelpStatusActive,
			Type:    model.HelpTypeArticle,
		})

		r := tester.HttpAdmin.Post("/v1/help", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		n := schema.Help{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		defer help.DeleteHelpById(n.Id)

		assert.Equal(t, tittle, n.Title)
		assert.Equal(t, content, n.Content)
		assert.Equal(t, model.HelpStatusActive, n.Status)
		assert.Equal(t, model.HelpTypeArticle, n.Type)
	}
}
