// Copyright 2019 Axetroy. All rights reserved. MIT license.
package report_schema

import (
	"github.com/axetroy/go-server/module/report/report_model"
)

type ReportPure struct {
	Id          string                    `json:"id"`
	Uid         string                    `json:"uid"`
	Title       string                    `json:"title"`
	Content     string                    `json:"content"`
	Type        report_model.ReportType   `json:"type"`
	Status      report_model.ReportStatus `json:"status"`
	Screenshots []string                  `json:"screenshots"`
	Locked      bool                      `json:"locked"`
}

type Report struct {
	ReportPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
