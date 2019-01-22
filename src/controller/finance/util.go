package finance

import "strings"

func GetTableName(currency string) string {
	return "finance_log_" + strings.ToLower(currency)
}
