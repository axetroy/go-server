// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package address_test

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/internal/app/user_server/controller/address"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAreaRouter(t *testing.T) {
	r := tester.HttpUser.Get("/v1/area", nil, nil)

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

	areaInfo := schema.Area{}

	assert.Nil(t, tester.Decode(res.Data, &areaInfo))
}

func TestFindAddress(t *testing.T) {
	{
		r, err := address.FindAddress("110101")

		assert.Nil(t, err)

		assert.Equal(t, &address.Area{
			Province: address.AreaStruct{Code: "110000", Name: "北京市"},
			City:     address.AreaStruct{Code: "110100", Name: "北京市"},
			Country:  address.AreaStruct{Code: "110101", Name: "东城区"},
			Addr:     "北京市东城区",
		}, r)
	}

	{
		r, err := address.FindAddress("450103")

		assert.Nil(t, err)

		assert.Equal(t, &address.Area{
			Province: address.AreaStruct{Code: "450000", Name: "广西壮族自治区"},
			City:     address.AreaStruct{Code: "450100", Name: "南宁市"},
			Country:  address.AreaStruct{Code: "450103", Name: "青秀区"},
			Addr:     "广西壮族自治区南宁市青秀区",
		}, r)
	}

	{
		r, err := address.FindAddress("123123")

		assert.Nil(t, r)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Sprintf("Invalid code: %d", 123123), err.Error())
	}
}
