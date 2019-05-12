package message_test

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/message"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
)

func TestDeleteByAdmin(t *testing.T) {
	var (
		adminUid string
	)

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

	context := controller.Context{
		Uid: adminUid,
	}

	var testMessage schema.Message

	// 创建一条个人信息
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := message.Create(context, message.CreateMessageParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testMessage = schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &testMessage))

		defer message.DeleteMessageById(testMessage.Id)

		assert.Equal(t, title, testMessage.Title)
		assert.Equal(t, content, testMessage.Content)
	}

	// 获取通知
	{
		n := model.Message{
			Id: testMessage.Id,
		}

		assert.Nil(t, service.Db.Model(&n).Where(&n).First(&n).Error)
	}

	// 删除通知
	{
		r := message.DeleteByAdmin(context, testMessage.Id)

		assert.Equal(t, "", r.Message)
		assert.Equal(t, schema.StatusSuccess, r.Status)
	}

	// 再次获取通知，这时候通知应该已经被删除了
	{
		n := model.Message{
			Id: testMessage.Id,
		}

		if err := service.Db.Model(&n).Where(&n).First(&n).Error; err != nil {
			assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
		} else {
			assert.Fail(t, "通知应该已被删除")
		}
	}
}

func TestDeleteByAdminRouter(t *testing.T) {
	var (
		adminToken  string
		messageInfo = schema.Message{}
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

		if _, er := util.ParseToken(util.TokenPrefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			adminToken = adminInfo.Token
		}
	}

	header := mocker.Header{
		"Authorization": util.TokenPrefix + " " + adminToken,
	}

	// 创建一条系统通知
	{
		var (
			title   = "test title"
			content = "test content"
		)

		body, _ := json.Marshal(&message.CreateMessageParams{
			Title:   title,
			Content: content,
		})

		r := tester.HttpAdmin.Post("/v1/message/create", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		assert.Nil(t, tester.Decode(res.Data, &messageInfo))

		defer message.DeleteMessageById(messageInfo.Id)

		assert.Equal(t, title, messageInfo.Title)
		assert.Equal(t, content, messageInfo.Content)
	}

	// 删除这条通知
	{
		r := tester.HttpAdmin.Delete("/v1/message/delete/"+messageInfo.Id, nil, &header)

		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		// 再查找这条记录，应该是空的

		n := model.Message{Id: messageInfo.Id}

		err := service.Db.Where(&n).First(&n).Error

		assert.NotNil(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
	}
}

func TestDeleteByUser(t *testing.T) {
	var (
		uid      string
		adminUid string
	)

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

	// 创建一个普通用户
	{
		var username = "test-delete-message"
		rand.Seed(10331198)
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

	context := controller.Context{
		Uid: uid,
	}

	var testMessage schema.Message

	// 创建一条个人信息
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := message.Create(controller.Context{
			Uid: adminUid,
		}, message.CreateMessageParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testMessage = schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &testMessage))

		defer message.DeleteMessageById(testMessage.Id)

		assert.Equal(t, title, testMessage.Title)
		assert.Equal(t, content, testMessage.Content)
	}

	// 获取通知
	{
		n := model.Message{
			Id: testMessage.Id,
		}

		assert.Nil(t, service.Db.Model(&n).Where(&n).First(&n).Error)
	}

	// 删除通知
	{
		r := message.DeleteByUser(context, testMessage.Id)

		assert.Equal(t, "", r.Message)
		assert.Equal(t, schema.StatusSuccess, r.Status)
	}

	// 再次获取通知，这时候通知应该已经被删除了
	{
		n := model.Message{
			Id: testMessage.Id,
		}

		if err := service.Db.Model(&n).Where(&n).First(&n).Error; err != nil {
			assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
		} else {
			assert.Fail(t, "通知应该已被删除")
		}
	}
}

func TestDeleteByUserRouter(t *testing.T) {
	var (
		messageInfo = schema.Message{}
	)

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": util.TokenPrefix + " " + adminInfo.Token,
	}

	fmt.Println(adminInfo)

	// 创建一条个人信息
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := message.Create(controller.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		messageInfo = schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &messageInfo))

		defer message.DeleteMessageById(messageInfo.Id)

		assert.Equal(t, title, messageInfo.Title)
		assert.Equal(t, content, messageInfo.Content)
	}

	// 删除这条通知
	{
		r := tester.HttpUser.Delete("/v1/message/delete/"+messageInfo.Id, nil, &header)

		res := schema.Response{}

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
			return
		}

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		fmt.Println(res)

		// 再查找这条记录，应该是空的

		n := model.Message{}

		err := service.Db.Where("id = ?", messageInfo.Id).First(&n).Error

		if !assert.NotNil(t, err) {
			return
		}
		if !assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error()) {
			return
		}
	}
}
