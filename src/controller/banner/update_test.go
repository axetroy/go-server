package banner_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/banner"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
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

func TestUpdate(t *testing.T) {

	var (
		adminUid   string
		bannerInfo = schema.Banner{}
	)

	// 1. 先登陆获取管理员的Token
	{
		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "admin",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := schema.AdminProfileWithToken{}

		assert.Nil(t, tester.Decode(r.Data, &adminInfo))

		if c, er := util.ParseToken(util.TokenPrefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			adminUid = c.Uid
		}
	}

	context := controller.Context{
		Uid: adminUid,
	}

	// 创建一个 Banner
	{
		var (
			image    = "test.png"
			href     = "https://example.com"
			platform = model.BannerPlatformApp
		)

		r := banner.Create(context, banner.CreateParams{
			Image:    image,
			Href:     href,
			Platform: platform,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &bannerInfo))

		defer banner.DeleteBannerById(bannerInfo.Id)

		assert.Equal(t, image, bannerInfo.Image)
		assert.Equal(t, href, bannerInfo.Href)
		assert.Equal(t, platform, bannerInfo.Platform)
	}

	// 更新这个刚添加的地址
	{

		var (
			newDescription = "new address"
		)

		r := banner.Update(context, bannerInfo.Id, banner.UpdateParams{
			Description: &newDescription,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &bannerInfo))

		assert.Equal(t, newDescription, *bannerInfo.Description)
	}

	{
		var (
			newDescription = "new new address"
			newHref        = "http://test.com"
		)

		r := banner.Update(context, bannerInfo.Id, banner.UpdateParams{
			Description: &newDescription,
			Href:        &newHref,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &bannerInfo))

		assert.Equal(t, newDescription, *bannerInfo.Description)
		assert.Equal(t, newHref, bannerInfo.Href)
	}
}

func TestUpdateRouter(t *testing.T) {
	var (
		tokenString string
		bannerInfo  = schema.Banner{}
	)

	// 1. 先登陆获取管理员的Token
	{
		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "admin",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := schema.AdminProfileWithToken{}

		assert.Nil(t, tester.Decode(r.Data, &adminInfo))

		if _, er := util.ParseToken(util.TokenPrefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			tokenString = adminInfo.Token
		}
	}

	header := mocker.Header{
		"Authorization": util.TokenPrefix + " " + tokenString,
	}

	// 创建一条 banner
	{

		body, _ := json.Marshal(&banner.CreateParams{
			Image:    "test.png",
			Href:     "http://example.org",
			Platform: model.BannerPlatformApp,
		})

		r := tester.Http.Post("/v1/admin/banner/create", body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		assert.Nil(t, tester.Decode(res.Data, &bannerInfo))

		defer banner.DeleteBannerById(bannerInfo.Id)
	}

	// 修改这条 banner
	{

		var (
			newImage       = "new.png"
			newDescription = "13333333333"
		)

		body, _ := json.Marshal(&banner.UpdateParams{
			Image:       &newImage,
			Description: &newDescription,
		})

		r := tester.Http.Put("/v1/admin/banner/update/"+bannerInfo.Id, body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		assert.Nil(t, tester.Decode(res.Data, &bannerInfo))

		assert.Equal(t, newImage, bannerInfo.Image)
		assert.Equal(t, newDescription, *bannerInfo.Description)

	}

}
