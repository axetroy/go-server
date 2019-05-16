package banner_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/banner"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/axetroy/go-server/src/util"
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

		defer banner.DeleteBannerById(n.Id)

		bannerId = n.Id

		assert.Equal(t, image, n.Image)
		assert.Equal(t, href, n.Href)
		assert.Equal(t, platform, n.Platform)
	}

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	// 删除这个刚添加的地址
	{
		r := banner.Delete(context, bannerId)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		bannerInfo := schema.Banner{}

		assert.Nil(t, tester.Decode(r.Data, &bannerInfo))

		assert.Equal(t, image, bannerInfo.Image)
		assert.Equal(t, href, bannerInfo.Href)
		assert.Equal(t, platform, bannerInfo.Platform)

		if err := service.Db.First(&model.Banner{
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

		defer banner.DeleteBannerById(n.Id)

		bannerId = n.Id

		assert.Equal(t, image, n.Image)
		assert.Equal(t, href, n.Href)
		assert.Equal(t, platform, n.Platform)
	}

	header := mocker.Header{
		"Authorization": util.TokenPrefix + " " + adminInfo.Token,
	}

	// 删除这条地址
	{

		r := tester.HttpAdmin.Delete("/v1/banner/b/"+bannerId, nil, &header)

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

		bannerInfo := schema.Banner{}

		assert.Nil(t, tester.Decode(res.Data, &bannerInfo))

		assert.Equal(t, image, bannerInfo.Image)
		assert.Equal(t, href, bannerInfo.Href)
		assert.Equal(t, platform, bannerInfo.Platform)

	}

}
