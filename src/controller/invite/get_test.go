package invite_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/invite"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	userInfo, err := tester.CreateUser()

	assert.Nil(t, err)

	defer auth.DeleteUserByUserName(userInfo.Username)

	// 获取一个不存在的邀请记录
	{
		r := invite.Get(controller.Context{
			Uid: userInfo.Id,
		}, "12313")

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.InviteNotExist.Error(), r.Message)
		assert.Nil(t, r.Data)
	}

	var inviteId1 string
	var inviteId2 string

	// 创建一条记录
	{
		tx := database.Db.Begin()

		v1 := model.InviteHistory{
			Inviter:       "123123",
			Invitee:       userInfo.Id, // 有一个跟测试账号相关的
			Status:        model.StatusInviteRegistered,
			RewardSettled: false,
		}

		v2 := model.InviteHistory{
			Inviter:       "123123", // 两个字段都测试账号不想关
			Invitee:       "123123",
			Status:        model.StatusInviteRegistered,
			RewardSettled: false,
		}

		if err := tx.Create(&v1).Error; err != nil {
			tx.Rollback()
			t.Error(err)
			return
		}

		if err := tx.Create(&v2).Error; err != nil {
			tx.Rollback()
			t.Error(err)
			return
		}

		tx.Commit()

		inviteId1 = v1.Id
		inviteId2 = v2.Id

		// 删除测试记录
		defer func() {
			invite.DeleteUserById(v1.Id)
			invite.DeleteUserById(v2.Id)
		}()
	}

	// 获取一个存在的
	{
		r := invite.Get(controller.Context{
			Uid: userInfo.Id,
		}, inviteId1)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		inviteInfo := schema.Invite{}

		assert.Nil(t, tester.Decode(r.Data, &inviteInfo))
	}

	// 获取一个跟我不相关的
	{
		r := invite.Get(controller.Context{
			Uid: userInfo.Id,
		}, inviteId2)

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.InviteNotExist.Error(), r.Message)
		assert.Nil(t, r.Data)
	}
}

func TestGetRouter(t *testing.T) {
	userInfo, err := tester.CreateUser()

	assert.Nil(t, err)

	defer auth.DeleteUserByUserName(userInfo.Username)

	// 获取一个不存在的邀请记录
	{
		r := invite.Get(controller.Context{
			Uid: userInfo.Id,
		}, "12313")

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.InviteNotExist.Error(), r.Message)
		assert.Nil(t, r.Data)
	}

	var inviteId string

	// 创建一条记录
	{
		tx := database.Db.Begin()

		v1 := model.InviteHistory{
			Inviter:       "123123",
			Invitee:       userInfo.Id, // 有一个跟测试账号相关的
			Status:        model.StatusInviteRegistered,
			RewardSettled: false,
		}

		v2 := model.InviteHistory{
			Inviter:       "123123", // 两个字段都测试账号不想关
			Invitee:       "123123",
			Status:        model.StatusInviteRegistered,
			RewardSettled: false,
		}

		if err := tx.Create(&v1).Error; err != nil {
			tx.Rollback()
			t.Error(err)
			return
		}

		if err := tx.Create(&v2).Error; err != nil {
			tx.Rollback()
			t.Error(err)
			return
		}

		tx.Commit()

		inviteId = v1.Id

		// 删除测试记录
		defer func() {
			invite.DeleteUserById(v1.Id)
			invite.DeleteUserById(v2.Id)
		}()
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	// 获取邀请方的记录
	{
		r := tester.HttpUser.Get("/v1/user/invite/i/"+inviteId, nil, &header)

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

		inviteDetail := schema.Invite{}

		assert.Nil(t, tester.Decode(res.Data, &inviteDetail))

		assert.Equal(t, "123123", inviteDetail.Inviter)
		assert.Equal(t, userInfo.Id, inviteDetail.Invitee)
	}

}
