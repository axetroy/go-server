package wallet

import (
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"strings"
	"time"
)

func GetTableName(currency string) string {
	return "wallet_" + strings.ToLower(currency)
}

func mapToSchema(model model.Wallet, d *schema.Wallet) {
	d.Id = model.Id
	d.Currency = model.Currency
	d.Balance = util.FloatToStr(model.Balance)
	d.Frozen = util.FloatToStr(model.Frozen)
	d.CreatedAt = model.CreatedAt.Format(time.RFC3339Nano)
	d.UpdatedAt = model.UpdatedAt.Format(time.RFC3339Nano)
}
