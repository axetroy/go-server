// Copyright 2019 Axetroy. All rights reserved. MIT license.
package wallet_test

import (
	"github.com/axetroy/go-server/core/controller/wallet"
	"github.com/axetroy/go-server/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTableName(t *testing.T) {
	assert.Equal(t, "wallet_cny", wallet.GetTableName(model.WalletCNY))
	assert.Equal(t, "wallet_usd", wallet.GetTableName(model.WalletUSD))
	assert.Equal(t, "wallet_coin", wallet.GetTableName(model.WalletCOIN))
}
