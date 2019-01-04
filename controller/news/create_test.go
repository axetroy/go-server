package news_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/controller/admin"
	"github.com/axetroy/go-server/controller/news"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/go-server/token"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	var (
		adminUid string
	)
	// 先登陆获取管理员的Token
	{
		// 登陆超级管理员-成功

		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "admin",
		})

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := admin.SignInResponse{}

		if err := tester.Decode(r.Data, &adminInfo); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, "admin", adminInfo.Username)
		assert.True(t, len(adminInfo.Token) > 0)

		if c, er := token.Parse(token.Prefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			adminUid = c.Uid
		}
	}

	// 创建一篇新闻
	{
		var (
			title    = "test"
			content  = "test"
			newsType = model.NewsType_News
		)

		r := news.Create(adminUid, news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    newsType,
			Tags:    []string{},
		})

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer func() {
			// TODO: 删除这篇文章
		}()

		n := news.News{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		assert.Equal(t, adminUid, n.Author)
		assert.Equal(t, title, n.Tittle)
		assert.Equal(t, content, n.Content)
		assert.Equal(t, newsType, n.Type)
		assert.Len(t, n.Tags, 0)
	}
}

func TestCreateRouter(t *testing.T) {

	var (
		adminUid   string
		adminToken string
	)
	// 先登陆获取管理员的Token
	{
		// 登陆超级管理员-成功

		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "admin",
		})

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := admin.SignInResponse{}

		if err := tester.Decode(r.Data, &adminInfo); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, "admin", adminInfo.Username)
		assert.True(t, len(adminInfo.Token) > 0)

		if c, er := token.Parse(token.Prefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			adminUid = c.Uid
			adminToken = adminInfo.Token
		}
	}

	// 登陆正确的管理员账号
	{
		var (
			title    = "test"
			content  = "test"
			newsType = model.NewsType_News
		)

		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminToken,
		}

		body, _ := json.Marshal(&news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    newsType,
			Tags:    []string{},
		})

		r := tester.Http.Post("/v1/admin/news/create", body, &header)
		res := response.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		defer func() {
			// TODO: 删除这篇文章
		}()

		n := news.News{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		assert.Equal(t, adminUid, n.Author)
		assert.Equal(t, title, n.Tittle)
		assert.Equal(t, content, n.Content)
		assert.Equal(t, newsType, n.Type)
		assert.Len(t, n.Tags, 0)
	}
}
