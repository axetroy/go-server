// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package banner_test

import (
	"encoding/json"
	banner2 "github.com/axetroy/go-server/internal/app/admin_server/controller/banner"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDelete(t *testing.T) {
	var (
		bannerId string
		image    = "test.png"
		href     = "https://example.com"
		platform = model.BannerPlatformApp
	)
	adminInfo, _ := tester.LoginAdmin()

	// 创建一个 Banner
	{
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

		bannerId = n.Id

		assert.Equal(t, image, n.Image)
		assert.Equal(t, href, n.Href)
		assert.Equal(t, platform, n.Platform)
	}

	context := helper.Context{
		Uid: adminInfo.Id,
	}

	// 删除这个刚添加的地址
	{
		r := banner2.Delete(context, bannerId)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		bannerInfo := schema.Banner{}

		assert.Nil(t, r.Decode(&bannerInfo))

		assert.Equal(t, image, bannerInfo.Image)
		assert.Equal(t, href, bannerInfo.Href)
		assert.Equal(t, platform, bannerInfo.Platform)

		if err := database.Db.First(&model.Banner{
			Id: bannerInfo.Id,
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
		bannerId string
		image    = "test.png"
		href     = "https://example.com"
		platform = model.BannerPlatformApp
	)
	adminInfo, _ := tester.LoginAdmin()

	// 创建一个 Banner
	{
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

		bannerId = n.Id

		assert.Equal(t, image, n.Image)
		assert.Equal(t, href, n.Href)
		assert.Equal(t, platform, n.Platform)
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	// 删除这条地址
	{

		r := tester.HttpAdmin.Delete("/v1/banner/"+bannerId, nil, &header)

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

		bannerInfo := schema.Banner{}

		assert.Nil(t, res.Decode(&bannerInfo))

		assert.Equal(t, image, bannerInfo.Image)
		assert.Equal(t, href, bannerInfo.Href)
		assert.Equal(t, platform, bannerInfo.Platform)

	}

}
