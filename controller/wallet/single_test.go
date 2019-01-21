package wallet_test

import (
	"github.com/axetroy/go-server/controller"
	"github.com/axetroy/go-server/controller/auth"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/controller/wallet"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestGetWallet(t *testing.T) {
	var (
		testUser user.Profile
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

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = user.Profile{}

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

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		walletInfo := wallet.Wallet{}

		assert.Nil(t, tester.Decode(r.Data, &walletInfo))

		assert.Equal(t, testUser.Id, walletInfo.Id)
		assert.Equal(t, float64(0), walletInfo.Balance)
		assert.Equal(t, float64(0), walletInfo.Frozen)
	}
}
