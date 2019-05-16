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

		assert.Nil(t, tester.Decode(r.Data, &addressInfo))

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

		assert.Nil(t, tester.Decode(r.Data, &addressInfo))

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

		r := tester.HttpUser.Put("/v1/user/address/a/"+addressInfo.Id, body, &header)

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

		assert.Equal(t, newName, addressInfo.Name)
		assert.Equal(t, newPhone, addressInfo.Phone)

	}

}
