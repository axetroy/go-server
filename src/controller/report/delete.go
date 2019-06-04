// Copyright 2019 Axetroy. All rights reserved. MIT license.
package report

import (
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/service/database"
)

func DeleteReportById(id string) {
	b := model.Report{}
	database.DeleteRowByTable(b.TableName(), "id", id)
}
