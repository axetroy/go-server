// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package address_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/user_server/controller/address"
	"github.com/axetroy/go-server/internal/app/user_server/controller/auth"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
)

func TestUpdate(t *testing.T) {

	var (
		username    = "tester-address-update-1"
		uid         string
		addressInfo = schema.Address{}
	)

	// 创建一个普通用户
	{
		rand.Seed(10331198)
		password := "123123"

		r := auth.SignUpWithUsername(auth.SignUpWithUsernameParams{
			Username: username,
			Password: password,
		})

		profile := schema.Profile{}

		assert.Nil(t, r.Decode(&profile))

		defer func() {
			tester.DeleteUserByUserName(username)
		}()

		uid = profile.Id
	}

	context := helper.Context{
		Uid: uid,
	}

	// 添加一个合法的地址
	{
		var (
			Name         = "test"
			Phone        = "13888888888"
			ProvinceCode = "11"
			CityCode     = "1101"
			AreaCode     = "110101"
			StreetCode   = "110101001"
			Address      = "中关村28号526"
		)

		r := address.Create(context, address.CreateAddressParams{
			Name:         Name,
			Phone:        Phone,
			ProvinceCode: ProvinceCode,
			CityCode:     CityCode,
			AreaCode:     AreaCode,
			StreetCode:   StreetCode,
			Address:      Address,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, r.Decode(&addressInfo))

		defer address.DeleteAddressById(addressInfo.Id)

	}

	// 更新这个刚添加的地址
	{

		var (
			newName = "new address"
		)

		r := address.Update(context, addressInfo.Id, address.UpdateParams{
			Name: &newName,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, r.Decode(&addressInfo))

		assert.Equal(t, newName, addressInfo.Name)
	}

	{
		var (
			newName  = "new new address"
			newPhone = "13333333333"
		)

		r := address.Update(context, addressInfo.Id, address.UpdateParams{
			Name:  &newName,
			Phone: &newPhone,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, r.Decode(&addressInfo))

		assert.Equal(t, newName, addressInfo.Name)
		assert.Equal(t, newPhone, addressInfo.Phone)
	}
}

func TestUpdateRouter(t *testing.T) {
	var (
		username    = "test-address-update-2"
		password    = "123123"
		tokenString string
		addressInfo = schema.Address{}
	)

	// 创建测试账号
	{
		if r := auth.SignUpWithUsername(auth.SignUpWithUsernameParams{
			Username: username,
			Password: password,
		}); r.Status != schema.StatusSuccess {
			t.Error(r.Message)
			return
		} else {
			userInfo := schema.Profile{}
			if err := r.Decode(&userInfo); err != nil {
				t.Error(err)
				return
			}
			defer func() {
				tester.DeleteUserByUserName(username)
			}()

			// 登陆获取Token
			if r := auth.SignIn(helper.Context{
				UserAgent: "test",
				Ip:        "0.0.0.0.0",
			}, auth.SignInParams{
				Account:  username,
				Password: password,
			}); r.Status != schema.StatusSuccess {
				t.Error(r.Message)
				return
			} else {
				userInfo := schema.ProfileWithToken{}
				if err := r.Decode(&userInfo); err != nil {
					t.Error(err)
					return
				}
				tokenString = userInfo.Token
			}
		}
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + tokenString,
	}

	// 创建一条收货地址
	{

		body, _ := json.Marshal(&address.CreateAddressParams{
			Name:         "张三",
			Phone:        "18888888888",
			ProvinceCode: "11",
			CityCode:     "1101",
			AreaCode:     "110101",
			StreetCode:   "110101001",
			Address:      "中关村28号526",
		})

		r := tester.HttpUser.Post("/v1/user/address", body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		assert.Nil(t, res.Decode(&addressInfo))

		defer address.DeleteAddressById(addressInfo.Id)
	}

	// 修改这条地址
	{

		var (
			newName  = "new address"
			newPhone = "13333333333"
		)

		body, _ := json.Marshal(&address.UpdateParams{
			Name:  &newName,
			Phone: &newPhone,
		})

		r := tester.HttpUser.Put("/v1/user/address/"+addressInfo.Id, body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		assert.Nil(t, res.Decode(&addressInfo))

		assert.Equal(t, newName, addressInfo.Name)
		assert.Equal(t, newPhone, addressInfo.Phone)

	}

}

func TestUpdateDefaultAddr(t *testing.T) {
	testUser, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(testUser.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + testUser.Token,
	}

	var addr1 schema.Address
	var addr2 schema.Address

	// 创建多个默认地址，理论上应该只有一个默认地址
	{
		isDefault := true
		body, _ := json.Marshal(&address.CreateAddressParams{
			Name:         "张三",
			Phone:        "18888888888",
			ProvinceCode: "11",
			CityCode:     "1101",
			AreaCode:     "110101",
			StreetCode:   "110101001",
			Address:      "中关村28号526",
			IsDefault:    &isDefault,
		})

		r := tester.HttpUser.Post("/v1/user/address", body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, "", res.Message)

		assert.Equal(t, schema.StatusSuccess, res.Status)

		addr1 = schema.Address{}

		assert.Nil(t, res.Decode(&addr1))

		defer address.DeleteAddressById(addr1.Id)

		assert.True(t, addr1.IsDefault)
	}

	// 创建多个默认地址，理论上应该只有一个默认地址
	{
		isDefault := true
		body, _ := json.Marshal(&address.CreateAddressParams{
			Name:         "张三",
			Phone:        "18888888888",
			ProvinceCode: "11",
			CityCode:     "1101",
			AreaCode:     "110101",
			StreetCode:   "110101001",
			Address:      "中关村28号526",
			IsDefault:    &isDefault,
		})

		r := tester.HttpUser.Post("/v1/user/address", body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, "", res.Message)

		assert.Equal(t, schema.StatusSuccess, res.Status)

		addr2 = schema.Address{}

		assert.Nil(t, res.Decode(&addr2))

		defer address.DeleteAddressById(addr2.Id)

		assert.True(t, addr2.IsDefault)
	}

	// 更新地址 1 为默认
	{
		isDefault := true
		body, _ := json.Marshal(&address.UpdateParams{
			IsDefault: &isDefault,
		})

		r := tester.HttpUser.Put("/v1/user/address/"+addr1.Id, body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		addressInfo := schema.Address{}

		assert.Nil(t, res.Decode(&addressInfo))

		assert.True(t, addressInfo.IsDefault)
	}

	// 获取详情 1
	{
		r := tester.HttpUser.Get("/v1/user/address/"+addr1.Id, nil, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		addressDetail := schema.Address{}

		assert.Nil(t, res.Decode(&addressDetail))

		assert.True(t, addressDetail.IsDefault)
	}

	// 获取详情 2
	{
		r := tester.HttpUser.Get("/v1/user/address/"+addr2.Id, nil, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		addressDetail := schema.Address{}

		assert.Nil(t, res.Decode(&addressDetail))

		assert.False(t, addressDetail.IsDefault)
	}

	// 获取列表
	{
		r := tester.HttpUser.Get("/v1/user/address", nil, &header)

		assert.Equal(t, http.StatusOK, r.Code)

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		addresses := make([]schema.Address, 0)

		assert.Nil(t, res.Decode(&addresses))

		for _, b := range addresses {
			switch b.Id {
			case addr1.Id:
				assert.True(t, b.IsDefault)
				break
			case addr2.Id:
				assert.False(t, b.IsDefault)
				break
			}
		}
	}
}
