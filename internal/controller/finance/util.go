// Copyright 2019 Axetroy. All rights reserved. MIT license.
package finance

import "strings"

func GetTableName(currency string) string {
	return "finance_log_" + strings.ToLower(currency)
}
