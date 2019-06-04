// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

import "github.com/axetroy/go-server/src/model"

type ReportPure struct {
	Id          string             `json:"id"`
	Uid         string             `json:"uid"`
	Title       string             `json:"title"`
	Content     string             `json:"content"`
	Type        model.ReportType   `json:"type"`
	Status      model.ReportStatus `json:"status"`
	Screenshots []string           `json:"screenshots"`
}

type Report struct {
	ReportPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
