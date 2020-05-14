// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package banner_test

import (
	"encoding/json"
	banner2 "github.com/axetroy/go-server/internal/app/admin_server/controller/banner"
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

func TestCreate(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	// 创建一个 Banner
	{
		var (
			image    = "test.png"
			href     = "https://example.com"
			platform = model.BannerPlatformApp
		)

		r := banner2.Create(helper.Context{
			Uid: adminInfo.Id,
		}, banner2.CreateParams{
			Image:    image,
			Href:     href,
			Platform: platform,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Banner{}

		assert.Nil(t, r.Decode(&n))

		defer banner2.DeleteBannerById(n.Id)

		assert.Equal(t, image, n.Image)
		assert.Equal(t, href, n.Href)
		assert.Equal(t, platform, n.Platform)
	}

	// 非管理员的uid去创建，应该报错
	{

		userInfo, _ := tester.CreateUser()

		defer tester.DeleteUserByUserName(userInfo.Username)

		var (
			image    = "test.png"
			href     = "https://example.com"
			platform = model.BannerPlatformApp
		)

		r := banner2.Create(helper.Context{
			Uid: userInfo.Id,
		}, banner2.CreateParams{
			Image:    image,
			Href:     href,
			Platform: platform,
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.AdminNotExist.Error(), r.Message)
	}
}

func TestCreateRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	// 创建 banner
	{
		var (
			image    = "test.png"
			href     = "https://example.com"
			platform = model.BannerPlatformApp
		)

		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		body, _ := json.Marshal(&banner2.CreateParams{
			Image:    image,
			Href:     href,
			Platform: platform,
		})

		r := tester.HttpAdmin.Post("/v1/banner", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		n := schema.Banner{}

		assert.Nil(t, res.Decode(&n))

		defer banner2.DeleteBannerById(n.Id)

		assert.Equal(t, image, n.Image)
		assert.Equal(t, href, n.Href)
		assert.Equal(t, platform, n.Platform)
	}
}
