// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package signature_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/controller"
	"github.com/axetroy/go-server/internal/controller/signature"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	r := signature.Encryption(controller.Context{}, "123")

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)
	assert.Equal(t, "aa154db7952aa6a2656fed90d0d88f2e87560a6ca7d7ed180ac76705fdc1639b", r.Data)
}

func TestCreateRouter(t *testing.T) {
	res := tester.HttpUser.Post("/v1/signature", []byte("123"), nil)
	r := schema.Response{}

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Nil(t, json.Unmarshal([]byte(res.Body.String()), &r))

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	assert.Equal(t, "aa154db7952aa6a2656fed90d0d88f2e87560a6ca7d7ed180ac76705fdc1639b", r.Data)

}
