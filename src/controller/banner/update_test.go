// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/banner"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUpdate(t *testing.T) {
	var (
		bannerInfo = schema.Banner{}
	)

	adminInfo, _ := tester.LoginAdmin()

	context := controller.Context{
		Uid: adminInfo.Id,
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
		bannerInfo = schema.Banner{}
	)

	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	// 创建一条 banner
	{

		body, _ := json.Marshal(&banner.CreateParams{
			Image:    "test.png",
			Href:     "http://example.org",
			Platform: model.BannerPlatformApp,
		})

		r := tester.HttpAdmin.Post("/v1/banner", body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes()), &res)) {
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

		r := tester.HttpAdmin.Put("/v1/banner/b/"+bannerInfo.Id, body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes()), &res)) {
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
