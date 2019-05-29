package news_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/news"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDelete(t *testing.T) {
	var (
		title   = "test"
		content = "test"
		newType = model.NewsType_Announcement
		tags    = []string{"test"}
		newsId  string
	)
	adminInfo, _ := tester.LoginAdmin()

	{

		r := news.Create(controller.Context{
			Uid: adminInfo.Id,
		}, news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    newType,
			Tags:    tags,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.News{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		newsId = n.Id

		defer news.DeleteNewsById(n.Id)

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
		assert.Equal(t, newType, n.Type)
		assert.Equal(t, tags, n.Tags)
	}

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	{
		r := news.Delete(context, newsId)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		newsInfo := schema.News{}

		assert.Nil(t, tester.Decode(r.Data, &newsInfo))

		assert.Equal(t, title, newsInfo.Title)
		assert.Equal(t, content, newsInfo.Content)
		assert.Equal(t, newType, newsInfo.Type)
		assert.Equal(t, tags, newsInfo.Tags)

		if err := database.Db.First(&model.News{
			Id: newsInfo.Id,
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
		title   = "test"
		content = "test"
		newType = model.NewsType_Announcement
		tags    = []string{"test"}
		newsId  string
	)
	adminInfo, _ := tester.LoginAdmin()

	{

		r := news.Create(controller.Context{
			Uid: adminInfo.Id,
		}, news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    newType,
			Tags:    tags,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.News{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		newsId = n.Id

		defer news.DeleteNewsById(n.Id)

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
		assert.Equal(t, newType, n.Type)
		assert.Equal(t, tags, n.Tags)
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	// 删除这条地址
	{

		r := tester.HttpAdmin.Delete("/v1/news/n/"+newsId, nil, &header)

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

		newsInfo := schema.News{}

		assert.Nil(t, tester.Decode(res.Data, &newsInfo))

		assert.Equal(t, title, newsInfo.Title)
		assert.Equal(t, content, newsInfo.Content)
		assert.Equal(t, newType, newsInfo.Type)
		assert.Equal(t, tags, newsInfo.Tags)
	}

}
