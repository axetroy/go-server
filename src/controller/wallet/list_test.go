// Copyright 2019 Axetroy. All rights reserved. MIT license.
package wallet_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/wallet"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWallets(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	r := wallet.GetWallets(controller.Context{Uid: userInfo.Id})

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	assert.Len(t, r.Data, len(model.WalletTableNames))

	list := make([]schema.Wallet, 0)
	assert.Nil(t, tester.Decode(r.Data, &list))

	assert.True(t, len(list) >= 1)

	for _, b := range list {
		assert.Equal(t, userInfo.Id, b.Id)
		assert.Equal(t, "0.00000000", b.Balance)
		assert.Equal(t, "0.00000000", b.Frozen)
		assert.IsType(t, "string", b.Currency)
		assert.IsType(t, "string", b.CreatedAt)
		assert.IsType(t, "string", b.UpdatedAt)
	}
}

func TestGetWalletsRouter(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	// 获取转账记录
	r := tester.HttpUser.Get("/v1/wallet", nil, &header)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes()), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	list := make([]schema.Wallet, 0)
	assert.Nil(t, tester.Decode(res.Data, &list))

	assert.True(t, len(list) >= 1)

	for _, b := range list {
		assert.Equal(t, userInfo.Id, b.Id)
		assert.Equal(t, "0.00000000", b.Balance)
		assert.Equal(t, "0.00000000", b.Frozen)
		assert.IsType(t, "string", b.Currency)
		assert.IsType(t, "string", b.CreatedAt)
		assert.IsType(t, "string", b.UpdatedAt)
	}
}
