package wallet_test

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/wallet"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestGetWallet(t *testing.T) {
	var (
		testUser schema.Profile
	)

	{
		// 1。 创建测试账号
		rand.Seed(111)
		username := "test-TestGetWallet"
		password := "123123"

		r := auth.SignUp(auth.SignUpParams{
			Username: &username,
			Password: password,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = schema.Profile{}

		if err := tester.Decode(r.Data, &testUser); err != nil {
			t.Error(err)
			return
		}

		defer func() {
			auth.DeleteUserByUserName(username)
		}()
	}

	{
		// 2. 获取获取钱包详情
		r := wallet.GetWallet(controller.Context{
			Uid: testUser.Id,
		}, model.WalletCNY)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		walletInfo := wallet.Wallet{}

		assert.Nil(t, tester.Decode(r.Data, &walletInfo))

		assert.Equal(t, testUser.Id, walletInfo.Id)
		assert.Equal(t, float64(0), walletInfo.Balance)
		assert.Equal(t, float64(0), walletInfo.Frozen)
	}
}
