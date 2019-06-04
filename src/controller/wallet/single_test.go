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
	"net/http"
	"testing"
)

func TestGetWallet(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	{
		// 2. 获取获取钱包详情
		r := wallet.GetWallet(controller.Context{
			Uid: userInfo.Id,
		}, model.WalletCNY)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		walletInfo := schema.Wallet{}

		assert.Nil(t, tester.Decode(r.Data, &walletInfo))

		assert.Equal(t, userInfo.Id, walletInfo.Id)
		assert.Equal(t, model.WalletCNY, walletInfo.Currency)
		assert.Equal(t, "0.00000000", walletInfo.Balance)
		assert.Equal(t, "0.00000000", walletInfo.Frozen)
		assert.NotEmpty(t, walletInfo.CreatedAt)
		assert.NotEmpty(t, walletInfo.UpdatedAt)
	}
}

func TestGetWalletRouter(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	// 获取详情
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + userInfo.Token,
		}

		r := tester.HttpUser.Get("/v1/wallet/w/cny", nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		n := schema.Wallet{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		assert.Equal(t, userInfo.Id, n.Id)
		assert.Equal(t, model.WalletCNY, n.Currency)
		assert.Equal(t, "0.00000000", n.Balance)
		assert.Equal(t, "0.00000000", n.Frozen)
		assert.Equal(t, "0.00000000", n.Frozen)
		assert.NotEmpty(t, n.CreatedAt)
		assert.NotEmpty(t, n.UpdatedAt)
	}
}
