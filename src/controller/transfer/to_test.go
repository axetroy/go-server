package transfer_test

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/transfer"
	"github.com/axetroy/go-server/src/controller/wallet"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTo(t *testing.T) {
	userFrom, _ := tester.CreateUser()
	userTo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userFrom.Username)
	defer auth.DeleteUserByUserName(userTo.Username)

	res1 := transfer.To(controller.Context{
		Uid: userFrom.Id,
	}, transfer.ToParams{
		Currency: "CNY",
		To:       userTo.Id,
		Amount:   "0.0001", // 转账失败，钱包没有余额
	})

	assert.Equal(t, exception.NotEnoughBalance.Error(), res1.Message)
	assert.Equal(t, schema.StatusFail, res1.Status)

	// 给账户充钱
	assert.Nil(t, service.Db.Table(wallet.GetTableName("CNY")).Where("id = ?", userFrom.Id).Update(model.Wallet{
		Balance:  100,
		Currency: model.WalletCNY,
	}).Error)

	res2 := transfer.To(controller.Context{
		Uid: userFrom.Id,
	}, transfer.ToParams{
		Currency: "CNY",
		To:       userTo.Id,
		Amount:   "20", // 转账 20
	})
	data := schema.TransferLog{}

	assert.Equal(t, "", res2.Message)
	assert.Equal(t, schema.StatusSuccess, res2.Status)
	assert.Nil(t, tester.Decode(res2.Data, &data))

	assert.Equal(t, userFrom.Id, data.From)
	assert.Equal(t, userTo.Id, data.To)
	assert.Equal(t, "20.00000000", data.Amount)

	// 检验账户金额是否正确
	r3 := wallet.GetWallet(controller.Context{Uid: userFrom.Id}, "CNY")
	fromUserWallet := schema.Wallet{}

	assert.Equal(t, "", r3.Message)
	assert.Equal(t, schema.StatusSuccess, r3.Status)
	assert.Nil(t, tester.Decode(r3.Data, &fromUserWallet))
	assert.Equal(t, "80.00000000", fromUserWallet.Balance)
	assert.Equal(t, "0.00000000", fromUserWallet.Frozen)

	r4 := wallet.GetWallet(controller.Context{Uid: userTo.Id}, "CNY")
	toUserWallet := schema.Wallet{}

	assert.Equal(t, "", r4.Message)
	assert.Equal(t, schema.StatusSuccess, r4.Status)
	assert.Nil(t, tester.Decode(r4.Data, &toUserWallet))
	assert.Equal(t, "20.00000000", toUserWallet.Balance)
	assert.Equal(t, "0.00000000", toUserWallet.Frozen)
}
