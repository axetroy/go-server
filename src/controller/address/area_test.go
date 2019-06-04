// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller/address"
	"github.com/axetroy/go-server/src/schema"
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

	if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
		return
	}

	if !assert.Equal(t, "", res.Message) {
		return
	}

	if !assert.Equal(t, schema.StatusSuccess, res.Status) {
		return
	}

	areaInfo := address.AreaListResponse{}

	assert.Nil(t, tester.Decode(res.Data, &areaInfo))
}
