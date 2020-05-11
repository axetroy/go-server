// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/controller"
	"github.com/axetroy/go-server/internal/controller/address"
	"github.com/axetroy/go-server/internal/controller/auth"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDelete(t *testing.T) {

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

	}

	// 删除这个刚添加的地址
	{
		r := address.Delete(context, addressInfo.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &addressInfo))

		assert.Equal(t, "test", addressInfo.Name)
		assert.Equal(t, "13888888888", addressInfo.Phone)

		addressInfo := model.Address{
			Id:  addressInfo.Id,
			Uid: context.Uid,
		}

		if err := database.Db.First(&addressInfo).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				assert.Fail(t, "数据被删除，应该不能再找到")
			}
		} else {
			assert.Fail(t, "数据被删除，应该不能再找到")
		}
	}

}

func TestDeleteRouter(t *testing.T) {
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

	// 创建一条收货地址
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

	// 删除这条地址
	{

		r := tester.HttpUser.Delete("/v1/user/address/a/"+addressInfo.Id, nil, &header)

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

		assert.Equal(t, "张三", addressInfo.Name)
		assert.Equal(t, "18888888888", addressInfo.Phone)

	}

}
