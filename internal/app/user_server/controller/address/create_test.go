// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package address_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/user_server/controller/address"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	testUser, err := tester.CreateUser()

	assert.Nil(t, err)

	defer tester.DeleteUserByUserName(testUser.Username)

	context := helper.Context{Uid: testUser.Id}

	// 添加一个失败的地址
	r := address.Create(context, address.CreateAddressParams{
		Name: "123",
	})

	assert.Equal(t, exception.InvalidParams.Code(), r.Status)

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

		addressInfo := schema.Address{}

		assert.Nil(t, r.Decode(&addressInfo))

		defer address.DeleteAddressById(addressInfo.Id)

		assert.Equal(t, Name, addressInfo.Name)
		assert.Equal(t, Phone, addressInfo.Phone)
		assert.Equal(t, ProvinceCode, addressInfo.ProvinceCode)
		assert.Equal(t, CityCode, addressInfo.CityCode)
		assert.Equal(t, AreaCode, addressInfo.AreaCode)
		// 之前没有添加地址的话，就是默认地址
		assert.Equal(t, true, addressInfo.IsDefault)
	}
}

func TestCreateRouter(t *testing.T) {
	testUser, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(testUser.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + testUser.Token,
	}

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

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

	assert.Equal(t, "", res.Message)

	assert.Equal(t, schema.StatusSuccess, res.Status)

	addressInfo := schema.Address{}

	assert.Nil(t, res.Decode(&addressInfo))

	defer address.DeleteAddressById(addressInfo.Id)
}

func TestCreateDefaultAddr(t *testing.T) {
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

		assert.False(t, addressDetail.IsDefault)
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

		assert.True(t, addressDetail.IsDefault)
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
				assert.False(t, b.IsDefault)
			case addr2.Id:
				assert.True(t, b.IsDefault)
			}
		}
	}
}
