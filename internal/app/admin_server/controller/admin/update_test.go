// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/admin"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUpdate(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	context := helper.Context{
		Uid: adminInfo.Id,
	}

	r := admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "123123",
		Password: "123123",
		Name:     "123123",
	}, false)

	assert.Equal(t, "", r.Message)
	assert.Equal(t, schema.StatusSuccess, r.Status)

	defer admin.DeleteAdminByAccount("123123")

	testAdminInfo := model.Admin{}

	assert.Nil(t, r.Decode(&testAdminInfo))

	status1 := model.AdminStatusInactivated

	res := admin.Update(context, testAdminInfo.Id, admin.UpdateParams{Status: &status1})

	assert.Equal(t, "", res.Message)
	assert.Equal(t, schema.StatusSuccess, res.Status)

	m := model.Admin{Id: testAdminInfo.Id}

	assert.Nil(t, database.Db.First(&m).Error)
	assert.Equal(t, status1, m.Status)
}

func TestUpdateRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	r1 := admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "123123",
		Password: "123123",
		Name:     "123123",
	}, false)

	assert.Equal(t, "", r1.Message)
	assert.Equal(t, schema.StatusSuccess, r1.Status)

	defer admin.DeleteAdminByAccount("123123")

	testAdminInfo := model.Admin{}

	assert.Nil(t, r1.Decode(&testAdminInfo))

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	newStatus := model.AdminStatusBanned

	body, _ := json.Marshal(&admin.UpdateParams{
		Status: &newStatus,
	})

	r := tester.HttpAdmin.Put("/v1/admin/"+testAdminInfo.Id, body, &header)

	res := schema.Response{}
	testProfile := schema.AdminProfile{}

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, "", res.Message)
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Nil(t, res.Decode(&testProfile))
	assert.Equal(t, newStatus, testProfile.Status)

}
