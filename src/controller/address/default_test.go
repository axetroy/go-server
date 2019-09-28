// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/address"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetDefault(t *testing.T) {

	var (
		addressInfo = schema.Address{}
	)

	userInfo, err := tester.CreateUser()

	if !assert.Nil(t, err) {
		return
	}

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := controller.Context{
		Uid: userInfo.Id,
	}

	// 还没有添加地址的话，获取默认地址应该会报错
	{
		r := address.GetDefault(context)

		assert.Equal(t, exception.AddressDefaultNotExist.Error(), r.Message)
		assert.Equal(t, schema.StatusFail, r.Status)
	}

	// 添加一个合法的地址
	{
		var (
			Name         = "test"
			Phone        = "13888888888"
			ProvinceCode = "110000"
			CityCode     = "110100"
			AreaCode     = "110101"
			Address      = "中关村28号526"
		)

		r := address.Create(context, address.CreateAddressParams{
			Name:         Name,
			Phone:        Phone,
			ProvinceCode: ProvinceCode,
			CityCode:     CityCode,
			AreaCode:     AreaCode,
			Address:      Address,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &addressInfo))

		defer address.DeleteAddressById(addressInfo.Id)

		assert.Equal(t, Name, addressInfo.Name)
		assert.Equal(t, Phone, addressInfo.Phone)
		assert.Equal(t, ProvinceCode, addressInfo.ProvinceCode)
		assert.Equal(t, CityCode, addressInfo.CityCode)
		assert.Equal(t, AreaCode, addressInfo.AreaCode)
		// 之前没有添加地址的话，就是默认地址
		assert.Equal(t, true, addressInfo.IsDefault)
	}

	// 获取默认地址
	{
		r := address.GetDefault(context)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		add := schema.Address{}

		assert.Nil(t, tester.Decode(r.Data, &add))

		assert.Equal(t, addressInfo.Name, add.Name)
		assert.Equal(t, addressInfo.Phone, add.Phone)
		assert.Equal(t, addressInfo.ProvinceCode, add.ProvinceCode)
		assert.Equal(t, addressInfo.CityCode, add.CityCode)
		assert.Equal(t, addressInfo.AreaCode, add.AreaCode)
		assert.Equal(t, true, add.IsDefault)
	}
}

func TestGetDefaultRouter(t *testing.T) {
	var (
		addressInfo = schema.Address{}
	)

	userInfo, err := tester.CreateUser()

	if !assert.Nil(t, err) {
		return
	}

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	// 创建一个地址
	{
		body, _ := json.Marshal(&address.CreateAddressParams{
			Name:         "张三",
			Phone:        "18888888888",
			ProvinceCode: "110000",
			CityCode:     "110100",
			AreaCode:     "110101",
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

		assert.Nil(t, tester.Decode(res.Data, &addressInfo))

		defer address.DeleteAddressById(addressInfo.Id)
	}

	// 获取默认地址
	{
		r := tester.HttpUser.Get("/v1/user/address/default", nil, &header)

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

		defaultAddress := schema.Address{}

		assert.Nil(t, tester.Decode(res.Data, &defaultAddress))

		assert.Equal(t, "张三", defaultAddress.Name)
		assert.Equal(t, "18888888888", defaultAddress.Phone)
	}
}
