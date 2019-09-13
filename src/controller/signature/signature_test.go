// Copyright 2019 Axetroy. All rights reserved. MIT license.
package signature_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/signature"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	r := signature.Encryption(controller.Context{}, "123")

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)
	assert.Equal(t, "38130d15223d921022b83673aaf26dde49bdb103e9471a1f52a801919364924a", r.Data)
}

func TestCreateRouter(t *testing.T) {
	res := tester.HttpUser.Post("/v1/signature", []byte("123"), nil)
	r := schema.Response{}

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Nil(t, json.Unmarshal([]byte(res.Body.String()), &r))

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	assert.Equal(t, "38130d15223d921022b83673aaf26dde49bdb103e9471a1f52a801919364924a", r.Data)

}
