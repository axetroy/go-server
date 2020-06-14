// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package migrate

import (
	"fmt"
	"time"
)

type Migrate interface {
	GetTableName() string                                  // 获取表名
	GetModel() interface{}                                 // 获取表名
	GetTimeInterval(now time.Time) time.Duration           // 时间间隔
	Next(nows ...time.Time) (shouldGoNext bool, err error) // 开始迁移, 可选的时间 nows 只是方便测试, 如果不传则使用 time.Now()
	Do() error                                             // 开始迁移，不断调用 Next()
}

var LoginLogMigrate Migrate
var CustomerMigrate Migrate

func init() {
	LoginLogMigrate = LoginLog{}
	CustomerMigrate = Customer{}
}

// 通过日期获取表名
func generateTableName(tableName string, date time.Time) string {
	year := fmt.Sprintf("%d", date.Year())
	month := fmt.Sprintf("%d", date.Month())
	if len(month) == 1 {
		month = "0" + month
	}

	newTableName := tableName + "_" + year + month

	return newTableName
}
