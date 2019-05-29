package banner_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/banner"
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

func TestGetBanner(t *testing.T) {
	{
		r := banner.GetBanner("123123")

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.BannerNotExist.Error(), r.Message)
	}

	{
		var (
			bannerId string
			image    = "test"
			href     = "test"
			platform = model.BannerPlatformApp
		)

		adminInfo, _ := tester.LoginAdmin()

		// 2. 先创建一篇新闻作为测试
		{

			r := banner.Create(controller.Context{
				Uid: adminInfo.Id,
			}, banner.CreateParams{
				Image:    image,
				Href:     href,
				Platform: platform,
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			n := schema.Banner{}

			assert.Nil(t, tester.Decode(r.Data, &n))

			bannerId = n.Id

			defer banner.DeleteBannerById(n.Id)
		}

		// 3. 获取文章公告
		{
			r := banner.GetBanner(bannerId)

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			bannerInfo := schema.Banner{}

			assert.Nil(t, tester.Decode(r.Data, &bannerInfo))

			assert.Equal(t, image, bannerInfo.Image)
			assert.Equal(t, href, bannerInfo.Href)
			assert.Equal(t, platform, bannerInfo.Platform)
		}
	}
}

func TestGetBannerRouter(t *testing.T) {
	var (
		bannerId string
		image    = "test"
		href     = "test"
		platform = model.BannerPlatformApp
	)

	adminInfo, _ := tester.LoginAdmin()

	// 先创建一篇新闻作为测试
	{

		r := banner.Create(controller.Context{
			Uid: adminInfo.Id,
		}, banner.CreateParams{
			Image:    image,
			Href:     href,
			Platform: platform,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Banner{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		bannerId = n.Id

		defer banner.DeleteBannerById(n.Id)
	}

	// 获取详情
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		r := tester.HttpAdmin.Get("/v1/banner/b/"+bannerId, nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		n := schema.Banner{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		assert.Equal(t, image, n.Image)
		assert.Equal(t, href, n.Href)
		assert.Equal(t, platform, n.Platform)
	}
}
