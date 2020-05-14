// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/util"
	"net/http"
	"testing"

	"github.com/axetroy/go-server/internal/app/admin_server/controller/admin"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePassword(t *testing.T) {
	r := admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "123123",
		Password: "123123",
		Name:     "123123",
	}, false)

	assert.Equal(t, "", r.Message)
	assert.Equal(t, schema.StatusSuccess, r.Status)

	defer admin.DeleteAdminByAccount("123123")

	testAdminInfo := schema.AdminProfile{}

	assert.Nil(t, r.Decode(&testAdminInfo))

	context := helper.Context{
		Uid: testAdminInfo.Id,
	}

	res := admin.UpdatePassword(context, admin.UpdatePasswordParams{
		OldPassword:     "123123",
		NewPassword:     "321321",
		ConfirmPassword: "321321",
	})

	assert.Equal(t, "", res.Message)
	assert.Equal(t, schema.StatusSuccess, res.Status)

	m := model.Admin{Id: testAdminInfo.Id}

	assert.Nil(t, database.Db.First(&m).Error)
	assert.Equal(t, util.GeneratePassword("321321"), m.Password)
}

func TestUpdatePasswordRouter(t *testing.T) {
	r1 := admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "123123",
		Password: "123123",
		Name:     "123123",
	}, false)

	assert.Equal(t, "", r1.Message)
	assert.Equal(t, schema.StatusSuccess, r1.Status)

	defer admin.DeleteAdminByAccount("123123")

	testAdminInfo := schema.AdminProfileWithToken{}

	assert.Nil(t, r1.Decode(&testAdminInfo))

	header := mocker.Header{
		"Authorization": token.Prefix + " " + testAdminInfo.Token,
	}

	body, _ := json.Marshal(&admin.UpdatePasswordParams{
		OldPassword:     "123123",
		NewPassword:     "321321",
		ConfirmPassword: "321321",
	})

	r := tester.HttpAdmin.Put("/v1/password", body, &header)

	res := schema.Response{}
	testProfile := schema.AdminProfile{}

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, "", res.Message)
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Nil(t, res.Decode(&testProfile))

	// 检查密码是否已被更改
}
