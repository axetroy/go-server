package invite_test

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/invite"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestGet(t *testing.T) {
	var (
		testUser schema.Profile
	)

	// 先创建一个测试用户
	{
		rand.Seed(111)
		username := "test-TestGet"
		password := "123123"

		r := auth.SignUp(auth.SignUpParams{
			Username: &username,
			Password: password,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = schema.Profile{}

		if err := tester.Decode(r.Data, &testUser); err != nil {
			t.Error(err)
			return
		}

		defer func() {
			auth.DeleteUserByUserName(username)
		}()
	}

	// 获取一个不存在的邀请记录
	{
		r := invite.Get(controller.Context{
			Uid: testUser.Id,
		}, "12313")

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.InviteNotExist.Error(), r.Message)
		assert.Nil(t, r.Data)
	}

	var inviteId1 string
	var inviteId2 string

	// 创建一条记录
	{
		tx := service.Db.Begin()

		v1 := model.InviteHistory{
			Inviter:       "123123",
			Invitee:       testUser.Id, // 有一个跟测试账号相关的
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
			Uid: testUser.Id,
		}, inviteId1)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		inviteInfo := schema.Invite{}

		assert.Nil(t, tester.Decode(r.Data, &inviteInfo))
	}

	// 获取一个跟我不相关的
	{
		r := invite.Get(controller.Context{
			Uid: testUser.Id,
		}, inviteId2)

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.NoPermission.Error(), r.Message)
		assert.Nil(t, r.Data)
	}
}
