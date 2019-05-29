package message_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/message"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetList(t *testing.T) {

	{
		var (
			data = make([]schema.Message, 0)
		)
		query := schema.Query{
			Limit: 20,
		}
		r := message.GetList(controller.Context{
			Uid: "123123",
		}, message.Query{
			Query: query,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &data))
		assert.Equal(t, query.Limit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.Equal(t, 0, r.Meta.Num)
		assert.Equal(t, int64(0), r.Meta.Total)
	}

	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(controller.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer message.DeleteMessageById(n.Id)
	}

	// 3. 获取列表
	{
		data := make([]schema.Message, 0)

		query := schema.Query{
			Limit: 20,
		}
		r := message.GetList(controller.Context{
			Uid: userInfo.Id,
		}, message.Query{
			Query: query,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &data))

		assert.Equal(t, query.Limit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.Equal(t, 1, r.Meta.Num)
		assert.Equal(t, int64(1), r.Meta.Total)

		assert.Len(t, data, 1)
	}
}

func TestGetListByAdmin(t *testing.T) {

	{
		var (
			data = make([]schema.Message, 0)
		)
		query := schema.Query{
			Limit: 20,
		}
		r := message.GetList(controller.Context{
			Uid: "123123",
		}, message.Query{
			Query: query,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &data))
		assert.Equal(t, query.Limit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.Equal(t, 0, r.Meta.Num)
		assert.Equal(t, int64(0), r.Meta.Total)
	}

	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(controller.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer message.DeleteMessageById(n.Id)
	}

	// 3. 获取列表
	{
		data := make([]schema.Message, 0)

		query := message.Query{
			Query: schema.Query{
				Limit: 20,
			},
		}
		r := message.GetListByAdmin(controller.Context{
			Uid: adminInfo.Id,
		}, message.QueryAdmin{
			Query: query,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &data))

		assert.Equal(t, query.Limit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.True(t, r.Meta.Num >= 1)
		assert.True(t, r.Meta.Total >= 1)
		assert.True(t, len(data) >= 1)
	}
}

func TestGetListRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(controller.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		//defer message.DeleteMessageById(n.Id)
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	{
		r := tester.HttpUser.Get("/v1/message", nil, &header)

		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		messages := make([]schema.Message, 0)

		assert.Nil(t, tester.Decode(res.Data, &messages))

		assert.True(t, len(messages) > 0)

		for _, b := range messages {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}

func TestGetListAdminRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(controller.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		//defer message.DeleteMessageById(n.Id)
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	{
		r := tester.HttpAdmin.Get("/v1/message", nil, &header)

		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		messages := make([]schema.Message, 0)

		assert.Nil(t, tester.Decode(res.Data, &messages))

		assert.True(t, len(messages) > 0)

		for _, b := range messages {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}
