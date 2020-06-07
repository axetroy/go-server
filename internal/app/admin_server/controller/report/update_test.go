// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package report_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/report"
	reportUser "github.com/axetroy/go-server/internal/app/user_server/controller/report"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUpdateByAdmin(t *testing.T) {
	var (
		title      = "title"
		content    = "content"
		reportType = model.ReportTypeBug
		reportInfo = schema.Report{}
	)
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	context := helper.Context{Uid: userInfo.Id}

	{
		r := reportUser.Create(context, reportUser.CreateParams{
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
		assert.Equal(t, model.ReportStatusPending, reportInfo.Status)
	}

	{
		locked := true
		r := report.UpdateByAdmin(helper.Context{Uid: adminInfo.Id}, reportInfo.Id, report.UpdateByAdminParams{
			Status: &model.ReportStatusResolve,
			Locked: &locked,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		reportInfo := schema.Report{}

		assert.Nil(t, r.Decode(&reportInfo))

		defer report.DeleteReportById(reportInfo.Id)

		assert.Equal(t, title, reportInfo.Title)
		assert.Equal(t, content, reportInfo.Content)
		assert.Equal(t, reportType, reportInfo.Type)
		assert.Equal(t, model.ReportStatusResolve, reportInfo.Status)
		assert.Equal(t, locked, reportInfo.Locked)
	}
}

func TestUpdateByAdminRouter(t *testing.T) {
	var (
		title      = "title"
		content    = "content"
		reportType = model.ReportTypeBug
		reportInfo = schema.Report{}
	)
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	context := helper.Context{Uid: userInfo.Id}

	{
		r := reportUser.Create(context, reportUser.CreateParams{
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
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		locked := true
		body, _ := json.Marshal(&report.UpdateByAdminParams{
			Status: &model.ReportStatusResolve,
			Locked: &locked,
		})

		res := tester.HttpAdmin.Put("/v1/report/"+reportInfo.Id, body, &header)
		r := schema.Response{}

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Nil(t, json.Unmarshal(res.Body.Bytes(), &r))

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		data := schema.Report{}

		assert.Nil(t, r.Decode(&data))

		assert.Equal(t, title, data.Title)
		assert.Equal(t, content, data.Content)
		assert.Equal(t, reportType, data.Type)
		assert.Equal(t, model.ReportStatusResolve, data.Status)
		assert.Equal(t, locked, data.Locked)
	}
}
