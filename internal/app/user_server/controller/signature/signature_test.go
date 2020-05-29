// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package signature_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/user_server/controller/signature"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	r := signature.Encryption(helper.Context{}, "123")

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	hash, err := util.Signature("123")
	assert.Nil(t, err)

	assert.Equal(t, hash, r.Data)
}

func TestCreateRouter(t *testing.T) {
	res := tester.HttpUser.Post("/v1/signature", []byte("123"), nil)
	r := schema.Response{}

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Nil(t, json.Unmarshal([]byte(res.Body.String()), &r))

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	hash, err := util.Signature("123")
	assert.Nil(t, err)

	assert.Equal(t, hash, r.Data)

}
