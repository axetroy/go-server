package address_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/address"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
)

func TestGetDetail(t *testing.T) {

	var (
		username    = "tester-normal-123"
		uid         string
		addressInfo = schema.Address{}
	)

	// 创建一个普通用户
	{
		rand.Seed(10331198)
		password := "123123"

		r := auth.SignUp(auth.SignUpParams{
			Username: &username,
			Password: password,
		})

		profile := schema.Profile{}

		assert.Nil(t, tester.Decode(r.Data, &profile))

		defer func() {
			auth.DeleteUserByUserName(username)
		}()

		uid = profile.Id
	}

	context := controller.Context{
		Uid: uid,
	}

	// 添加一个合法的地址
	{
		var (
			Name         = "test"
			Phone        = "13888888888"
			ProvinceCode = "100000"
			CityCode     = "101000"
			AreaCode     = "101010"
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

	// 获取地址详情
	{
		r := address.GetDetail(context, addressInfo.Id)

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

func TestGetDetailRouter(t *testing.T) {
	var (
		username    = "test-create-address"
		password    = "123123"
		tokenString string
		addressInfo = schema.Address{}
	)

	// 创建测试账号
	{
		if r := auth.SignUp(auth.SignUpParams{
			Username: &username,
			Password: password,
		}); r.Status != schema.StatusSuccess {
			t.Error(r.Message)
			return
		} else {
			userInfo := schema.Profile{}
			if err := tester.Decode(r.Data, &userInfo); err != nil {
				t.Error(err)
				return
			}
			defer func() {
				auth.DeleteUserByUserName(username)
			}()

			// 登陆获取Token
			if r := auth.SignIn(controller.Context{
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
				if err := tester.Decode(r.Data, &userInfo); err != nil {
					t.Error(err)
					return
				}
				tokenString = userInfo.Token
			}
		}

	}
	header := mocker.Header{
		"Authorization": util.TokenPrefix + " " + tokenString,
	}

	// 创建一个地址
	{
		body, _ := json.Marshal(&address.CreateAddressParams{
			Name:         "张三",
			Phone:        "18888888888",
			ProvinceCode: "100000",
			CityCode:     "101000",
			AreaCode:     "101010",
			Address:      "中关村28号526",
		})

		r := tester.Http.Post("/v1/user/address/create", body, &header)

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

		assert.Nil(t, tester.Decode(res.Data, &addressInfo))

		defer address.DeleteAddressById(addressInfo.Id)
	}

	// 获取详情
	{
		r := tester.Http.Get("/v1/user/address/detail/"+addressInfo.Id, nil, &header)

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

		addressDetail := schema.Address{}

		assert.Nil(t, tester.Decode(res.Data, &addressDetail))

		assert.Equal(t, "张三", addressDetail.Name)
		assert.Equal(t, "18888888888", addressDetail.Phone)
	}
}
