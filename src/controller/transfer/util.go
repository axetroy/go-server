package transfer

import (
	"strings"
)

// 获取转账表名
func GetTransferTableName(currency string) string {
	return "transfer_log_" + strings.ToLower(currency)
}
