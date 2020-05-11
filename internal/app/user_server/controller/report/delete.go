// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package report

import (
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
)

func DeleteReportById(id string) {
	b := model.Report{}
	database.DeleteRowByTable(b.TableName(), "id", id)
}
