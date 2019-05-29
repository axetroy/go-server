package news_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/news"
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

	defer auth.DeleteUserByUserName(userInfo.Username)

	// 创建一个公告
	{
		var (
			title   = "test"
			content = "test"
		)

		r := news.Create(controller.Context{
			Uid: adminInfo.Id,
		}, news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    model.NewsType_News,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := model.News{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer news.DeleteNewsById(n.Id)

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
	}

	// 非管理员的uid去创建，应该报错
	{
		var (
			title    = "test"
			content  = "test"
			newsType = model.NewsType_News
		)

		r := news.Create(controller.Context{
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
			Type:    model.NewsType_News,
		})

		r := tester.HttpAdmin.Post("/v1/news", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		n := schema.News{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		defer news.DeleteNewsById(n.Id)

		assert.Equal(t, adminInfo.Id, n.Author)
		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
	}
}
