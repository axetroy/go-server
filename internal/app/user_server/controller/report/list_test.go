// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package report_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/user_server/controller/report"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	context := helper.Context{
		Uid: userInfo.Id,
	}

	{
		var (
			title      = "title"
			content    = "content"
			reportType = model.ReportTypeBug
			reportInfo = schema.Report{}
		)

		r := report.Create(context, report.CreateParams{
			Title:   title,
			Content: content,
			Type:    reportType,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, r.Decode(&reportInfo))

		defer report.DeleteReportById(reportInfo.Id)

		assert.Equal(t, title, reportInfo.Title)
		assert.Equal(t, content, reportInfo.Content)
		assert.Equal(t, reportType, reportInfo.Type)
	}

	// 获取列表
	{
		r := report.GetList(context, report.Query{})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		reports := make([]schema.Report, 0)

		assert.Nil(t, r.Decode(&reports))

		assert.Equal(t, schema.DefaultLimit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.IsType(t, 1, r.Meta.Num)
		assert.IsType(t, int64(1), r.Meta.Total)

		assert.True(t, len(reports) >= 1)

		for _, b := range reports {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, model.ReportTypeBug, b.Type)
			assert.IsType(t, model.ReportStatusPending, b.Status)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}

func TestGetListRouter(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	context := helper.Context{
		Uid: userInfo.Id,
	}

	{
		var (
			title      = "title"
			content    = "content"
			reportType = model.ReportTypeBug
			reportInfo = schema.Report{}
		)

		r := report.Create(context, report.CreateParams{
			Title:   title,
			Content: content,
			Type:    reportType,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, r.Decode(&reportInfo))

		defer report.DeleteReportById(reportInfo.Id)

		assert.Equal(t, title, reportInfo.Title)
		assert.Equal(t, content, reportInfo.Content)
		assert.Equal(t, reportType, reportInfo.Type)
	}

	{
		r := tester.HttpUser.Get("/v1/report", nil, &header)

		res := schema.Response{}

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		reports := make([]schema.Report, 0)

		assert.Nil(t, res.Decode(&reports))

		for _, b := range reports {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, model.ReportTypeBug, b.Type)
			assert.IsType(t, model.ReportStatusPending, b.Status)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}
