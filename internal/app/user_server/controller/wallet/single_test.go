// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package wallet_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/user_server/controller/wallet"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetWallet(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	{
		// 2. 获取获取钱包详情
		r := wallet.GetWallet(helper.Context{
			Uid: userInfo.Id,
		}, model.WalletCNY)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		walletInfo := schema.Wallet{}

		assert.Nil(t, r.Decode(&walletInfo))

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

	defer tester.DeleteUserByUserName(userInfo.Username)

	// 获取详情
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + userInfo.Token,
		}

		r := tester.HttpUser.Get("/v1/wallet/cny", nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		n := schema.Wallet{}

		assert.Nil(t, res.Decode(&n))

		assert.Equal(t, userInfo.Id, n.Id)
		assert.Equal(t, model.WalletCNY, n.Currency)
		assert.Equal(t, "0.00000000", n.Balance)
		assert.Equal(t, "0.00000000", n.Frozen)
		assert.Equal(t, "0.00000000", n.Frozen)
		assert.NotEmpty(t, n.CreatedAt)
		assert.NotEmpty(t, n.UpdatedAt)
	}
}
