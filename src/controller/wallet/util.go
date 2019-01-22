package wallet

import (
	"strings"
)

func GetTableName(currency string) string {
	return "wallet_" + strings.ToLower(currency)
}
