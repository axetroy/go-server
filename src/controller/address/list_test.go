// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/address"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetList(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := controller.Context{
		Uid: userInfo.Id,
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

		addressInfo := schema.Address{}

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

	// 获取地址列表
	{
		r := address.GetAddressListByUser(context, address.Query{})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		addresses := make([]schema.Address, 0)

		assert.Nil(t, tester.Decode(r.Data, &addresses))

		assert.Equal(t, schema.DefaultLimit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.Equal(t, 1, r.Meta.Num)
		assert.Equal(t, int64(1), r.Meta.Total)

		assert.Len(t, addresses, 1)

		//assert.Len(t, 1, len(addresses))

		firstAddress := addresses[0]

		assert.Equal(t, "test", firstAddress.Name)
		assert.Equal(t, "13888888888", firstAddress.Phone)
		assert.Equal(t, "110000", firstAddress.ProvinceCode)
		assert.Equal(t, "110100", firstAddress.CityCode)
		assert.Equal(t, "110101", firstAddress.AreaCode)
		assert.Equal(t, "中关村28号526", firstAddress.Address)
		assert.Equal(t, true, firstAddress.IsDefault)
	}
}

func TestGetListRouter(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := controller.Context{
		Uid: userInfo.Id,
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
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

		addressInfo := schema.Address{}

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

	{
		r := tester.HttpUser.Get("/v1/user/address", nil, &header)

		assert.Equal(t, http.StatusOK, r.Code)

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes()), &res)) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		addresses := make([]schema.Address, 0)

		assert.Nil(t, tester.Decode(res.Data, &addresses))

		for _, b := range addresses {
			assert.IsType(t, "string", b.Phone)
			assert.IsType(t, "string", b.Phone)
			assert.IsType(t, "string", b.ProvinceCode)
			assert.IsType(t, "string", b.CityCode)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}
