package news_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/news"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
)

func init() {
	// 确保超级管理员存在
	admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "admin",
		Password: "admin",
		Name:     "admin",
	}, true)
}

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

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := schema.AdminProfileWithToken{}

		if err := tester.Decode(r.Data, &adminInfo); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, "admin", adminInfo.Username)
		assert.True(t, len(adminInfo.Token) > 0)

		if c, er := util.ParseToken(util.TokenPrefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			adminUid = c.Uid
		}
	}

	// 创建一个公告
	{
		var (
			title   = "test"
			content = "test"
		)

		r := news.Create(controller.Context{
			Uid: adminUid,
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
		// 创建一个普通用户
		var (
			username = "tester-normal"
			uid      string
		)

		{
			rand.Seed(10331)
			password := "123123"

			r := auth.SignUp(auth.SignUpParams{
				Username: &username,
				Password: password,
			})

			profile := schema.Profile{}

			assert.Nil(t, tester.Decode(r.Data, &profile))

			defer auth.DeleteUserByUserName(username)

			uid = profile.Id
		}

		var (
			title    = "test"
			content  = "test"
			newsType = model.NewsType_News
		)

		r := news.Create(controller.Context{
			Uid: uid,
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

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := schema.AdminProfileWithToken{}

		if err := tester.Decode(r.Data, &adminInfo); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, "admin", adminInfo.Username)
		assert.True(t, len(adminInfo.Token) > 0)

		if c, er := util.ParseToken(util.TokenPrefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			adminUid = c.Uid
			adminToken = adminInfo.Token
		}
	}

	// 登陆正确的管理员账号
	{
		var (
			title   = "test"
			content = "test"
		)

		header := mocker.Header{
			"Authorization": util.TokenPrefix + " " + adminToken,
		}

		body, _ := json.Marshal(&news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    model.NewsType_News,
		})

		r := tester.Http.Post("/v1/admin/news/create", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		n := schema.News{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		defer news.DeleteNewsById(n.Id)

		assert.Equal(t, adminUid, n.Author)
		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
	}
}
