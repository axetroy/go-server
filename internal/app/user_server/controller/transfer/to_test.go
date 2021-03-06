// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package transfer_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/user_server/controller/transfer"
	"github.com/axetroy/go-server/internal/app/user_server/controller/user"
	"github.com/axetroy/go-server/internal/app/user_server/controller/wallet"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestTo(t *testing.T) {
	userFrom, _ := tester.CreateUser()
	userTo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userFrom.Username)
	defer tester.DeleteUserByUserName(userTo.Username)

	input1 := transfer.ToParams{
		Currency: "CNY",
		To:       userTo.Id,
		Amount:   "0.0001", // 转账失败，钱包没有余额
	}

	b1, err := json.Marshal(input1)

	assert.Nil(t, err)

	signature1, err := util.Signature(string(b1))

	assert.Nil(t, err)

	res1 := transfer.To(helper.Context{
		Uid: userFrom.Id,
	}, input1, signature1)

	assert.Equal(t, exception.NotEnoughBalance.Error(), res1.Message)
	assert.Equal(t, exception.NotEnoughBalance.Code(), res1.Status)

	// 给账户充钱
	assert.Nil(t, database.Db.Table(wallet.GetTableName("CNY")).Where("id = ?", userFrom.Id).Update(model.Wallet{
		Balance:  100,
		Currency: model.WalletCNY,
	}).Error)

	input2 := transfer.ToParams{
		Currency: "CNY",
		To:       userTo.Id,
		Amount:   "20", // 转账 20
	}

	b2, err := json.Marshal(input2)

	assert.Nil(t, err)

	signature2, err := util.Signature(string(b2))

	assert.Nil(t, err)

	res2 := transfer.To(helper.Context{
		Uid: userFrom.Id,
	}, input2, signature2)
	data := schema.TransferLog{}

	assert.Equal(t, "", res2.Message)
	assert.Equal(t, schema.StatusSuccess, res2.Status)
	assert.Nil(t, res2.Decode(&data))

	assert.Equal(t, userFrom.Id, data.From)
	assert.Equal(t, userTo.Id, data.To)
	assert.Equal(t, "20.00000000", data.Amount)

	// 检验账户金额是否正确
	r3 := wallet.GetWallet(helper.Context{Uid: userFrom.Id}, "CNY")
	fromUserWallet := schema.Wallet{}

	assert.Equal(t, "", r3.Message)
	assert.Equal(t, schema.StatusSuccess, r3.Status)
	assert.Nil(t, r3.Decode(&fromUserWallet))
	assert.Equal(t, "80.00000000", fromUserWallet.Balance)
	assert.Equal(t, "0.00000000", fromUserWallet.Frozen)

	r4 := wallet.GetWallet(helper.Context{Uid: userTo.Id}, "CNY")
	toUserWallet := schema.Wallet{}

	assert.Equal(t, "", r4.Message)
	assert.Equal(t, schema.StatusSuccess, r4.Status)
	assert.Nil(t, r4.Decode(&toUserWallet))
	assert.Equal(t, "20.00000000", toUserWallet.Balance)
	assert.Equal(t, "0.00000000", toUserWallet.Frozen)

	// Invalid Signature
	{
		input := transfer.ToParams{
			Currency: "CNY",
			To:       userTo.Id,
			Amount:   "0.0001", // 转账失败，钱包没有余额
		}

		res := transfer.To(helper.Context{
			Uid: userFrom.Id,
		}, input, "Invalid signature")

		assert.Equal(t, exception.InvalidSignature.Error(), res.Message)
		assert.Equal(t, exception.InvalidSignature.Code(), res.Status)
	}
}

func TestToRouter(t *testing.T) {
	userFrom, _ := tester.CreateUser()
	userTo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userFrom.Username)
	defer tester.DeleteUserByUserName(userTo.Username)

	// 设置用户的交易密码
	rr := user.SetPayPassword(helper.Context{Uid: userFrom.Id}, user.SetPayPasswordParams{
		Password:        "123123",
		PasswordConfirm: "123123",
	})

	assert.Equal(t, "", rr.Message)
	assert.Equal(t, schema.StatusSuccess, rr.Status)

	// 给账户充钱
	assert.Nil(t, database.Db.Table(wallet.GetTableName("CNY")).Where("id = ?", userFrom.Id).Update(model.Wallet{
		Balance:  100,
		Currency: model.WalletCNY,
	}).Error)

	// 转账
	{
		header := mocker.Header{
			"Authorization":  token.Prefix + " " + userFrom.Token,
			"X-Pay-Password": "123123",
		}

		body, _ := json.Marshal(&transfer.ToParams{
			Currency: "CNY",
			To:       userTo.Id,
			Amount:   "20",
		})

		signature, err := util.Signature(string(body))

		header[middleware.SignatureHeader] = signature

		assert.Nil(t, err)

		r := tester.HttpUser.Post("/v1/transfer", body, &header)

		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, "", res.Message)
		assert.Equal(t, schema.StatusSuccess, res.Status)

		// 检验账户金额是否正确
		r3 := wallet.GetWallet(helper.Context{Uid: userFrom.Id}, "CNY")
		fromUserWallet := schema.Wallet{}

		assert.Equal(t, "", r3.Message)
		assert.Equal(t, schema.StatusSuccess, r3.Status)
		assert.Nil(t, r3.Decode(&fromUserWallet))
		assert.Equal(t, "80.00000000", fromUserWallet.Balance)
		assert.Equal(t, "0.00000000", fromUserWallet.Frozen)

		r4 := wallet.GetWallet(helper.Context{Uid: userTo.Id}, "CNY")
		toUserWallet := schema.Wallet{}

		assert.Equal(t, "", r4.Message)
		assert.Equal(t, schema.StatusSuccess, r4.Status)
		assert.Nil(t, r4.Decode(&toUserWallet))
		assert.Equal(t, "20.00000000", toUserWallet.Balance)
		assert.Equal(t, "0.00000000", toUserWallet.Frozen)
	}
}
