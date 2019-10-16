// Copyright 2019 Axetroy. All rights reserved. MIT license.
package transfer

import (
	"fmt"
	"github.com/axetroy/go-server/core/model"
	"github.com/axetroy/go-server/core/util"
	"reflect"
	"strings"
)

// 获取转账表名
func GetTransferTableName(currency string) string {
	return "transfer_log_" + strings.ToLower(currency)
}

type QueryParams struct {
	Id       *string               `json:"id"`       // 转账ID
	Currency *string               `json:"currency"` // 转账币种
	From     *string               `json:"from"`     // 汇款人
	To       *string               `json:"to"`       // 收款人
	Amount   *string               `json:"amount"`   // 转账数量
	Status   *model.TransferStatus `json:"status"`   // 转账状态
}

func GenerateTransferLogSQL(filter QueryParams, limit int, count bool) string {
	suffix := `("deleted_at" IS NULL OR "deleted_at"='0001-01-01 00:00:00')`

	filterArray := make([]string, 0)

	{
		t := reflect.TypeOf(filter)
		v := reflect.ValueOf(filter)

		for k := 0; k < t.NumField(); k++ {
			key := t.Field(k).Tag.Get("json")
			value := v.Field(k).Interface()

			if key == "" {
				continue
			}

			if !v.Field(k).IsValid() {
				continue
			}

			if util.IsNil(value) {
				continue
			}

			// 如果是指针的话
			if util.IsPoint(value) {
				// 获取指针对应的值
				value = reflect.ValueOf(value).Elem()
			} else {
				continue
			}

			filterArray = append(filterArray, fmt.Sprintf(`"%s"='%v'`, key, value))
		}
	}

	filterStr := "WHERE " + strings.Join(filterArray[:], " AND ")

	if len(filterArray) != 0 {
		filterStr = filterStr + " AND"
	}

	SQLs := make([]string, 0)

	selected := "*"

	if count {
		selected = "COUNT(*)"
	}

	for _, tableName := range model.TransferTableNames {
		sql := fmt.Sprintf(`SELECT %s FROM "%s" %s %s`, selected, tableName, filterStr, suffix)
		SQLs = append(SQLs, sql)
	}

	sql := fmt.Sprintf("%s LIMIT %d", strings.Join(SQLs[:], " UNION "), limit)
	return sql
}
