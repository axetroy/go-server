package tester

import (
	"errors"

	"github.com/axetroy/go-server/internal/app/admin_server/controller/admin"
	"github.com/axetroy/go-server/internal/app/user_server/controller/auth"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
)

// 创建一个测试用户
func CreateUser() (profile schema.ProfileWithToken, err error) {
	var (
		username  = "test-" + util.RandomString(6)
		password  = "123123"
		ip        = "0.0.0.0"
		userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3788.1 Safari/537.36"
	)

	// 创建用户
	if r := auth.SignUpWithUsername(auth.SignUpWithUsernameParams{
		Username: username,
		Password: password,
	}); r.Status != schema.StatusSuccess {
		err = errors.New(r.Message)
		return
	}

	// 登陆获取 token
	r := auth.SignIn(helper.Context{
		UserAgent: userAgent,
		Ip:        ip,
	}, auth.SignInParams{
		Account:  username,
		Password: password,
	})

	if r.Status != schema.StatusSuccess {
		err = errors.New(r.Message)
		return
	}

	if err = r.Decode(&profile); err != nil {
		return
	}

	return
}

// 登陆超级管理员
func LoginAdmin() (profile schema.AdminProfileWithToken, err error) {
	r := admin.Login(admin.SignInParams{
		Username: "admin",
		Password: "123456",
	})

	if r.Status != schema.StatusSuccess {
		err = errors.New(r.Message)
		return
	}

	if err = r.Decode(&profile); err != nil {
		return
	}

	return
}

// 删除用户
func DeleteUserByUserName(username string) {
	database.DeleteRowByTable("user", "username", username)
}

// 删除用户
func DeleteUserByUid(uid string) {
	database.DeleteRowByTable("user", "id", uid)
	database.DeleteRowByTable("wechat_open_id", "uid", uid)
}
