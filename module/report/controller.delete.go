// Copyright 2019 Axetroy. All rights reserved. MIT license.
package report

import (
	"github.com/axetroy/go-server/module/report/report_model"
	"github.com/axetroy/go-server/service/database"
)

func DeleteReportById(id string) {
	b := report_model.Report{}
	database.DeleteRowByTable(b.TableName(), "id", id)
}
