// Copyright 2019 Axetroy. All rights reserved. MIT license.
package transfer_test

import (
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/transfer"
	"github.com/axetroy/go-server/module/transfer/transfer_schema"
	"github.com/axetroy/go-server/module/wallet"
	"github.com/axetroy/go-server/module/wallet/wallet_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDetail(t *testing.T) {
	var log transfer_schema.TransferLog
	userFrom, _ := tester.CreateUser()
	userTo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userFrom.Username)
	defer auth.DeleteUserByUserName(userTo.Username)

	// 给账户充钱
	{
		assert.Nil(t, database.Db.Table(wallet.GetTableName("CNY")).Where("id = ?", userFrom.Id).Update(wallet_model.Wallet{
			Balance:  100,
			Currency: wallet_model.WalletCNY,
		}).Error)
	}

	// 转账一次
	{
		res2 := transfer.To(schema.Context{
			Uid: userFrom.Id,
		}, transfer.ToParams{
			Currency: "CNY",
			To:       userTo.Id,
			Amount:   "20", // 转账 20
		})

		assert.Equal(t, "", res2.Message)
		assert.Equal(t, schema.StatusSuccess, res2.Status)
		assert.Nil(t, tester.Decode(res2.Data, &log))
	}

	{
		r := transfer.GetDetail(schema.Context{
			Uid: userFrom.Id,
		}, log.Id)

		detail := transfer_schema.TransferLog{}

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &detail))

		assert.Equal(t, log.Id, detail.Id)
		assert.Equal(t, log.From, detail.From)
		assert.Equal(t, log.To, detail.To)
		assert.Equal(t, log.Amount, detail.Amount)
		assert.Equal(t, log.Status, detail.Status)
	}
}
