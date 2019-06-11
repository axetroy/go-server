// Copyright 2019 Axetroy. All rights reserved. MIT license.
package report_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/report"
	"github.com/axetroy/go-server/module/report/report_model"
	"github.com/axetroy/go-server/module/report/report_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUpdate(t *testing.T) {
	var (
		title      = "title"
		content    = "content"
		reportType = report_model.ReportTypeBug
		reportInfo = report_schema.Report{}
	)
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := schema.Context{Uid: userInfo.Id}

	{
		r := report.Create(context, report.CreateParams{
			Title:   title,
			Content: content,
			Type:    reportType,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &reportInfo))

		defer report.DeleteReportById(reportInfo.Id)

		assert.Equal(t, title, reportInfo.Title)
		assert.Equal(t, content, reportInfo.Content)
		assert.Equal(t, reportType, reportInfo.Type)
		assert.Equal(t, report_model.ReportStatusPending, reportInfo.Status)
	}

	{
		r := report.Update(context, reportInfo.Id, report.UpdateParams{
			Status: &report_model.ReportStatusResolve,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		reportInfo := report_schema.Report{}

		assert.Nil(t, tester.Decode(r.Data, &reportInfo))

		defer report.DeleteReportById(reportInfo.Id)

		assert.Equal(t, title, reportInfo.Title)
		assert.Equal(t, content, reportInfo.Content)
		assert.Equal(t, reportType, reportInfo.Type)
		assert.Equal(t, report_model.ReportStatusResolve, reportInfo.Status)
	}
}

func TestUpdateRouter(t *testing.T) {
	var (
		title      = "title"
		content    = "content"
		reportType = report_model.ReportTypeBug
		reportInfo = report_schema.Report{}
	)
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := schema.Context{Uid: userInfo.Id}

	{
		r := report.Create(context, report.CreateParams{
			Title:   title,
			Content: content,
			Type:    reportType,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &reportInfo))

		defer report.DeleteReportById(reportInfo.Id)

		assert.Equal(t, title, reportInfo.Title)
		assert.Equal(t, content, reportInfo.Content)
		assert.Equal(t, reportType, reportInfo.Type)
	}

	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + userInfo.Token,
		}

		body, _ := json.Marshal(&report.UpdateParams{
			Status: &report_model.ReportStatusResolve,
		})

		res := tester.HttpUser.Put("/v1/report/r/"+reportInfo.Id, body, &header)
		r := schema.Response{}

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Nil(t, json.Unmarshal([]byte(res.Body.String()), &r))

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		data := report_schema.Report{}

		assert.Nil(t, tester.Decode(r.Data, &data))

		assert.Equal(t, title, data.Title)
		assert.Equal(t, content, data.Content)
		assert.Equal(t, reportType, data.Type)
		assert.Equal(t, report_model.ReportStatusResolve, data.Status)
	}
}

func TestUpdateByAdmin(t *testing.T) {
	var (
		title      = "title"
		content    = "content"
		reportType = report_model.ReportTypeBug
		reportInfo = report_schema.Report{}
	)
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	ctx := schema.Context{Uid: userInfo.Id}

	{
		r := report.Create(ctx, report.CreateParams{
			Title:   title,
			Content: content,
			Type:    reportType,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &reportInfo))

		defer report.DeleteReportById(reportInfo.Id)

		assert.Equal(t, title, reportInfo.Title)
		assert.Equal(t, content, reportInfo.Content)
		assert.Equal(t, reportType, reportInfo.Type)
		assert.Equal(t, report_model.ReportStatusPending, reportInfo.Status)
	}

	{
		locked := true
		r := report.UpdateByAdmin(schema.Context{Uid: adminInfo.Id}, reportInfo.Id, report.UpdateByAdminParams{
			UpdateParams: report.UpdateParams{
				Status: &report_model.ReportStatusResolve,
			},
			Locked: &locked,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		reportInfo := report_schema.Report{}

		assert.Nil(t, tester.Decode(r.Data, &reportInfo))

		defer report.DeleteReportById(reportInfo.Id)

		assert.Equal(t, title, reportInfo.Title)
		assert.Equal(t, content, reportInfo.Content)
		assert.Equal(t, reportType, reportInfo.Type)
		assert.Equal(t, report_model.ReportStatusResolve, reportInfo.Status)
		assert.Equal(t, locked, reportInfo.Locked)
	}
}

func TestUpdateByAdminRouter(t *testing.T) {
	var (
		title      = "title"
		content    = "content"
		reportType = report_model.ReportTypeBug
		reportInfo = report_schema.Report{}
	)
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := schema.Context{Uid: userInfo.Id}

	{
		r := report.Create(context, report.CreateParams{
			Title:   title,
			Content: content,
			Type:    reportType,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &reportInfo))

		defer report.DeleteReportById(reportInfo.Id)

		assert.Equal(t, title, reportInfo.Title)
		assert.Equal(t, content, reportInfo.Content)
		assert.Equal(t, reportType, reportInfo.Type)
	}

	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		locked := true
		body, _ := json.Marshal(&report.UpdateByAdminParams{
			UpdateParams: report.UpdateParams{
				Status: &report_model.ReportStatusResolve,
			},
			Locked: &locked,
		})

		res := tester.HttpAdmin.Put("/v1/report/r/"+reportInfo.Id, body, &header)
		r := schema.Response{}

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Nil(t, json.Unmarshal([]byte(res.Body.String()), &r))

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		data := report_schema.Report{}

		assert.Nil(t, tester.Decode(r.Data, &data))

		assert.Equal(t, title, data.Title)
		assert.Equal(t, content, data.Content)
		assert.Equal(t, reportType, data.Type)
		assert.Equal(t, report_model.ReportStatusResolve, data.Status)
		assert.Equal(t, locked, data.Locked)
	}
}
