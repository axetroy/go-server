// Copyright 2019 Axetroy. All rights reserved. MIT license.
package wallet_test

import (
	"github.com/axetroy/go-server/module/wallet"
	"github.com/axetroy/go-server/module/wallet/wallet_model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTableName(t *testing.T) {
	assert.Equal(t, "wallet_cny", wallet.GetTableName(wallet_model.WalletCNY))
	assert.Equal(t, "wallet_usd", wallet.GetTableName(wallet_model.WalletUSD))
	assert.Equal(t, "wallet_coin", wallet.GetTableName(wallet_model.WalletCOIN))
}
