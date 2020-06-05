// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package banner_test

import (
	"encoding/json"
	bannerAdmin "github.com/axetroy/go-server/internal/app/admin_server/controller/banner"
	"github.com/axetroy/go-server/internal/app/user_server/controller/banner"
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

func TestGetBanner(t *testing.T) {
	{
		r := banner.GetBanner("123123")

		assert.Equal(t, exception.BannerNotExist.Code(), r.Status)
		assert.Equal(t, exception.BannerNotExist.Error(), r.Message)
	}

	{
		var (
			bannerId string
			image    = "https://example/test.png"
			href     = "https://example.com"
			platform = model.BannerPlatformApp
		)

		adminInfo, _ := tester.LoginAdmin()

		// 2. 先创建一篇新闻作为测试
		{

			r := bannerAdmin.Create(helper.Context{
				Uid: adminInfo.Id,
			}, bannerAdmin.CreateParams{
				Image:    image,
				Href:     href,
				Platform: platform,
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			n := schema.Banner{}

			assert.Nil(t, r.Decode(&n))

			bannerId = n.Id

			defer bannerAdmin.DeleteBannerById(n.Id)
		}

		// 3. 获取文章公告
		{
			r := banner.GetBanner(bannerId)

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			bannerInfo := schema.Banner{}

			assert.Nil(t, r.Decode(&bannerInfo))

			assert.Equal(t, image, bannerInfo.Image)
			assert.Equal(t, href, bannerInfo.Href)
			assert.Equal(t, platform, bannerInfo.Platform)
		}
	}
}

func TestGetBannerRouter(t *testing.T) {
	var (
		bannerId string
		image    = "https://example/test.png"
		href     = "https://example.com"
		platform = model.BannerPlatformApp
	)

	adminInfo, _ := tester.LoginAdmin()

	// 先创建一篇新闻作为测试
	{

		r := bannerAdmin.Create(helper.Context{
			Uid: adminInfo.Id,
		}, bannerAdmin.CreateParams{
			Image:    image,
			Href:     href,
			Platform: platform,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Banner{}

		assert.Nil(t, r.Decode(&n))

		bannerId = n.Id

		defer bannerAdmin.DeleteBannerById(n.Id)
	}

	// 获取详情
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		r := tester.HttpAdmin.Get("/v1/banner/"+bannerId, nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		n := schema.Banner{}

		assert.Nil(t, res.Decode(&n))

		assert.Equal(t, image, n.Image)
		assert.Equal(t, href, n.Href)
		assert.Equal(t, platform, n.Platform)
	}
}
