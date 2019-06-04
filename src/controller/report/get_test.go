// Copyright 2019 Axetroy. All rights reserved. MIT license.
package report_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/report"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetReportByUser(t *testing.T) {
	var (
		title      = "title"
		content    = "content"
		reportType = model.ReportTypeBug
		reportInfo = schema.Report{}
	)
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := controller.Context{Uid: userInfo.Id}

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
		r := report.GetReportByUser(context, reportInfo.Id)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		data := schema.Report{}
		assert.Nil(t, tester.Decode(r.Data, &data))

		assert.Equal(t, title, data.Title)
		assert.Equal(t, content, data.Content)
		assert.Equal(t, reportType, data.Type)
	}
}

func TestGetReportRouter(t *testing.T) {
	var (
		title      = "title"
		content    = "content"
		reportType = model.ReportTypeBug
		reportInfo = schema.Report{}
	)
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := controller.Context{Uid: userInfo.Id}

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

		body, _ := json.Marshal(&report.CreateParams{
			Title:   title,
			Content: content,
			Type:    reportType,
		})

		res := tester.HttpUser.Get("/v1/report/r/"+reportInfo.Id, body, &header)
		r := schema.Response{}

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Nil(t, json.Unmarshal([]byte(res.Body.String()), &r))

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		data := schema.Report{}

		assert.Nil(t, tester.Decode(r.Data, &data))

		assert.Equal(t, title, data.Title)
		assert.Equal(t, content, data.Content)
		assert.Equal(t, reportType, data.Type)
	}
}
