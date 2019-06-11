// Copyright 2019 Axetroy. All rights reserved. MIT license.
package wallet_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/wallet"
	"github.com/axetroy/go-server/module/wallet/wallet_model"
	"github.com/axetroy/go-server/module/wallet/wallet_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWallets(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	r := wallet.GetWallets(schema.Context{Uid: userInfo.Id})

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	assert.Len(t, r.Data, len(wallet_model.WalletTableNames))

	list := make([]wallet_schema.Wallet, 0)
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

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	list := make([]wallet_schema.Wallet, 0)
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
