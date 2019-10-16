// Copyright 2019 Axetroy. All rights reserved. MIT license.
package transfer_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/core/controller"
	"github.com/axetroy/go-server/core/controller/auth"
	"github.com/axetroy/go-server/core/controller/transfer"
	"github.com/axetroy/go-server/core/controller/wallet"
	"github.com/axetroy/go-server/core/model"
	"github.com/axetroy/go-server/core/schema"
	"github.com/axetroy/go-server/core/service/database"
	"github.com/axetroy/go-server/core/service/token"
	"github.com/axetroy/go-server/core/util"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHistory(t *testing.T) {
	userFrom, _ := tester.CreateUser()
	userTo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userFrom.Username)
	defer auth.DeleteUserByUserName(userTo.Username)

	// 给账户充钱
	assert.Nil(t, database.Db.Table(wallet.GetTableName("CNY")).Where("id = ?", userFrom.Id).Update(model.Wallet{
		Balance:  100,
		Currency: model.WalletCNY,
	}).Error)

	// 创建一条转账记录
	input := transfer.ToParams{
		Currency: "CNY",
		To:       userTo.Id,
		Amount:   "20", // 转账 20
	}

	b, err := json.Marshal(input)

	assert.Nil(t, err)

	signature, err := util.Signature(string(b))

	assert.Nil(t, err)

	res2 := transfer.To(controller.Context{
		Uid: userFrom.Id,
	}, input, signature)

	assert.Equal(t, "", res2.Message)
	assert.Equal(t, schema.StatusSuccess, res2.Status)

	// 获取转账记录
	r := transfer.GetHistory(controller.Context{
		Uid: userFrom.Id,
	}, transfer.Query{})

	assert.Equal(t, "", r.Message)
	assert.Equal(t, schema.StatusSuccess, r.Status)

	list := make([]schema.TransferLog, 0)
	assert.Nil(t, tester.Decode(r.Data, &list))

	assert.Equal(t, schema.DefaultLimit, r.Meta.Limit)
	assert.Equal(t, schema.DefaultPage, r.Meta.Page)
	assert.IsType(t, 1, r.Meta.Num)
	assert.IsType(t, int64(1), r.Meta.Total)

	assert.True(t, len(list) >= 1)

	for _, b := range list {
		assert.Equal(t, userFrom.Id, b.From)
		assert.Equal(t, userTo.Id, b.To)
		assert.Equal(t, "CNY", b.Currency)
		assert.Equal(t, "20.00000000", b.Amount)
		assert.IsType(t, "string", b.CreatedAt)
		assert.IsType(t, "string", b.UpdatedAt)
	}
}

func TestGetHistoryRouter(t *testing.T) {
	userFrom, _ := tester.CreateUser()
	userTo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userFrom.Username)
	defer auth.DeleteUserByUserName(userTo.Username)

	// 给账户充钱
	assert.Nil(t, database.Db.Table(wallet.GetTableName("CNY")).Where("id = ?", userFrom.Id).Update(model.Wallet{
		Balance:  100,
		Currency: model.WalletCNY,
	}).Error)

	// 创建一条转账记录
	input := transfer.ToParams{
		Currency: "CNY",
		To:       userTo.Id,
		Amount:   "20", // 转账 20
	}

	b, err := json.Marshal(input)

	assert.Nil(t, err)

	signature, err := util.Signature(string(b))

	assert.Nil(t, err)

	res2 := transfer.To(controller.Context{
		Uid: userFrom.Id,
	}, input, signature)

	assert.Equal(t, "", res2.Message)
	assert.Equal(t, schema.StatusSuccess, res2.Status)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userFrom.Token,
	}

	// 获取转账记录
	r := tester.HttpUser.Get("/v1/transfer", nil, &header)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	list := make([]schema.TransferLog, 0)

	assert.Nil(t, tester.Decode(res.Data, &list))

	assert.True(t, len(list) >= 1)

	for _, b := range list {
		assert.Equal(t, userFrom.Id, b.From)
		assert.Equal(t, userTo.Id, b.To)
		assert.Equal(t, "CNY", b.Currency)
		assert.Equal(t, "20.00000000", b.Amount)
		assert.IsType(t, "string", b.CreatedAt)
		assert.IsType(t, "string", b.UpdatedAt)
	}
}
